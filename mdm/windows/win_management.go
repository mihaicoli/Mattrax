package windows

import (
	"context"
	"net/http"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/mdm/windows/syncml"
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
		// 	log.Debug().Str("protocol_udid", cmd.Header.SourceURI).Msg("Invalid device authenitcation missing TLS authentication certificate")
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

func ManagementHandler(ctx context.Context, srv *mattrax.Server, cmd syncml.Envelop, res syncml.Response, device db.Device) {
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
			// r.SetRaw("Add", "./Vendor/MSFT/NodeCache/MattraxMDM/Nodes/"+node+"/NodeURI", "", "", payload.Uri)
			// r.SetRaw("Add", "./Vendor/MSFT/NodeCache/MattraxMDM/Nodes/"+node+"/ExpectedValue", "", "", payload.Value)
		}

		if err := srv.DB.NewDeviceCacheNode(ctx, db.NewDeviceCacheNodeParams{
			DeviceID:  device.ID,
			PayloadID: payload.ID,
			// CacheID:   , // TODO: Nodecache location
		}); err != nil {
			log.Error().Err(err).Msg("Error updating device cache node")
			res.SetStatus(syncml.StatusCommandFailed)
			return
		}
		// TODO: NodeCache + Atomics
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
			PayloadID: payload.ID,
		}); err != nil {
			log.Error().Err(err).Msg("Error updating device cache node")
			res.SetStatus(syncml.StatusCommandFailed)
			return
		}
	}
}
