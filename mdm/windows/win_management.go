package windows

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/pkg/syncml"
	"github.com/rs/zerolog/log"
)

// Manage talks to the device periodically to update configuration and get device information
func Manage(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd, errored := syncml.Read(r, w)
		if errored {
			return
		}
		var res = syncml.NewResponse(cmd)

		// TODO: Certificate Authentication Not Working Correctly So Has Been Disabled
		// if len(r.TLS.PeerCertificates) == 0 {
		// 	log.Debug().Str("protocol_udid", cmd.Header.SourceURI).Msg("Invalid device authentication missing TLS authentication certificate")
		// 	var res = syncml.NewResponse(cmd)
		// 	res.Status(syncml.StatusUnauthorized)
		// 	res.Respond(w)
		// 	return
		// } else if r.TLS.PeerCertificates[0].Subject.CommonName != cmd.Header.SourceURI || r.TLS.PeerCertificates[0].Subject.OrganizationalUnit[0] != "WinMDM" {
		// 	log.Debug().Int("cert-count", len(r.TLS.PeerCertificates)).Str("CN", r.TLS.PeerCertificates[0].Subject.CommonName).Str("OU", r.TLS.PeerCertificates[0].Subject.OrganizationalUnit[0]).Str("protocol_udid", cmd.Header.SourceURI).Msg("Invalid device authentication")
		// 	var res = syncml.NewResponse(cmd)
		// 	res.Status(syncml.StatusUnauthorized)
		// 	res.Respond(w)
		// 	return
		// }

		// TODO: Check Request Data (Correct Destinations, etc)

		device, err := srv.DB.GetDeviceByUDID(r.Context(), cmd.Header.SourceURI)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving managed device")
			res.SetStatus(syncml.StatusCommandFailed)
			res.Respond(w)
			return
		}

		ManagementHandler(r.Context(), srv, cmd, res, device)

		if err := srv.DB.DeviceCheckinStatus(r.Context(), db.DeviceCheckinStatusParams{
			ID:             device.ID,
			LastseenStatus: res.FinalStatus(),
		}); err != nil {
			log.Error().Err(err).Msg("Error storing checkin status")
			res.SetStatus(syncml.StatusCommandFailed)
			res.Respond(w)
			return
		}

		res.Respond(w)
	}
}

// ManagementHandler handles deploying configuration and handling its response from the device
func ManagementHandler(ctx context.Context, srv *mattrax.Server, cmd syncml.Message, res syncml.Response, device db.Device) {
	if device.State == db.DeviceStateDeploying {
		if err := srv.DB.SetDeviceState(ctx, db.SetDeviceStateParams{
			ID:    device.ID,
			State: db.DeviceStateManaged,
		}); err != nil {
			log.Error().Err(err).Msg("Error retrieving managed device")
			res.SetStatus(syncml.StatusCommandFailed)
			return
		}
	} else if device.State != db.DeviceStateManaged && device.State != db.DeviceStateMissing {
		log.Error().Int32("id", device.ID).Str("state", string(device.State)).Msg("Device is unable to be managed due to unsupported state")
		res.SetStatus(syncml.StatusForbidden)
		return
	}

	// TODO: Make this look nicer
	for _, command := range cmd.Body.Commands {
		var final bool
		switch command.XMLName.Local {
		case "Alert":
			if command.Data == "1201" {
				continue
			} else if command.Data == "1224" {
				// TODO: ADD login status
				if command.Source != nil && command.Source.URI != "" && command.Meta != nil && command.Meta.Type == "come.microsoft.mdm.win32csp_install" {
					fmt.Println("MSI Install Status", command.Meta.Format, command.Meta.Mark, command.Data)
					// TODO: Handle This
				}
				// TODO: MBES Apps
				continue
			} else if command.Data == "1226" {
				// TODO: Check for body element
				if command.Body[0].Meta.Type == "com.microsoft:mdm.unenrollment.userrequest" {
					if err := srv.DB.DeviceUserUnenrollment(ctx, device.ID); err != nil {
						log.Error().Int32("id", device.ID).Err(err).Msg("Device is unable to be updated during unenrollment. This request will NEVER be sent again from the device!")
						res.SetStatus(syncml.StatusCommandFailed)
						return
					}
					log.Info().Int32("id", device.ID).Str("trigger", "user_enroll").Msg("Device unenrolled")
				}
				// TODO: ADD login status
				continue
			} else {
				fmt.Println("Unknown Alert Type:", command.Data)
			}
			break
		case "Results":
			fmt.Println(command)
			for _, command := range command.Body {
				if err := srv.DB.UpdateDeviceInventoryNode(ctx, db.UpdateDeviceInventoryNodeParams{
					DeviceID: device.ID,
					Uri:      command.Source.URI,
					Format:   command.Meta.Format,
					Value:    command.Data,
				}); err != nil {
					log.Error().Int32("id", device.ID).Str("uri", command.Source.URI).Err(err).Msg("Unable to update device inventory node")
					res.SetStatus(syncml.StatusCommandFailed)
					return
				}
			}
		case "Final":
			final = true
			break
		default:
			fmt.Println("Unsupported Command:", command.XMLName.Local)
			// 	UnsupportedCommand(command, state)
		}
		if final == true {
			break
		}
	}

	// TODO: NodeCache global version check like CSP defines server should do

	// TODO: Work out data that the inventory needs about the device then ask for it and NodeCache!
	// TODO: This includes DM CSP Versions

	// TODO: Parse Request Data and Store in Inventory + Detect and handle User Unenroll + Add/Replace switch

	payloadsAwaitingDeploy, err := srv.DB.GetDevicesPayloadsAwaitingDeployment(ctx, device.ID) // TODO: Replace SQL type because its a required arg
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving devices policies that are awaiting deploy")
		return
	}

	for _, payload := range payloadsAwaitingDeploy {
		if payload.Exec {
			res.Set("Add", payload.Uri, "", "", "")
			res.Set("Exec", payload.Uri, payload.Type, payload.Format, payload.Value)
		} else {
			res.Set("Add", payload.Uri, payload.Type, payload.Format, payload.Value)

			// TODO: NodeCache
			// r.SetRaw("Add", "./Vendor/MSFT/NodeCache/" + ProviderID + "/Nodes/"+node+"/NodeURI", "", "", payload.Uri)
			// r.SetRaw("Add", "./Vendor/MSFT/NodeCache/" + ProviderID + "/Nodes/"+node+"/ExpectedValue", "", "", payload.Value)
		}

		var nodecacheNode int32
		if nodecacheNode, err = srv.DB.NewDeviceCacheNode(ctx, db.NewDeviceCacheNodeParams{
			DeviceID:  device.ID,
			PayloadID: sql.NullInt32{payload.ID, true},
		}); err != nil {
			log.Error().Err(err).Msg("Error updating device cache node")
			res.SetStatus(syncml.StatusCommandFailed)
			return
		}

		fmt.Println("NodeCache Destined" + strconv.Itoa(int(nodecacheNode))) // TODO: NodeCache + Use Atomics for it
		// TODO: Method for checking if nodecache was set/if it successed?
	}

	detachedPayloads, err := srv.DB.GetDevicesDetachedPayloads(ctx, device.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving devices detached payloads")
		return
	}
	for _, payload := range detachedPayloads {
		res.Set("Delete", payload.Uri, "", "", "")
		// TODO: Remove NodeCache nodes

		if err := srv.DB.DeleteDeviceCacheNode(ctx, db.DeleteDeviceCacheNodeParams{
			DeviceID:  device.ID,
			PayloadID: sql.NullInt32{payload.ID, true},
		}); err != nil {
			log.Error().Err(err).Msg("Error updating device cache node")
			res.SetStatus(syncml.StatusCommandFailed)
			return
		}
	}
}
