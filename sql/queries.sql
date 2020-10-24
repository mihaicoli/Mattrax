-- DO NOT RUN THIS FILE. It is used along with sqlc to generate type safe Go from SQL

-- name: GetUsers :many
-- Exposed via API
SELECT upn, fullname FROM users LIMIT 100;

-- name: GetUser :one
-- Exposed via API
SELECT upn, fullname, azuread_oid FROM users WHERE upn = $1 LIMIT 1;

-- name: GetUserForLogin :one
SELECT fullname, password, mfa_token FROM users WHERE upn = $1 LIMIT 1;

-- name: NewAzureADUser :one
INSERT INTO users(upn, fullname, azuread_oid) VALUES($1, $2, $3) RETURNING upn, fullname, azuread_oid; -- TODO: Insert or Update

-- TODO: Merge all NewDevice functions to single query

-- name: NewDevice :one
INSERT INTO devices(udid, state, enrollment_type, name, hw_dev_id, operating_system, azure_did, enrolled_by) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;

-- name: NewDeviceReplacingExisting :exec
UPDATE devices SET state=$2, enrollment_type=$3, name=$4, hw_dev_id=$5, operating_system=$6, azure_did=$7, nodecache_version='', lastseen=NOW(), lastseen_status=0, enrolled_at=NOW(), enrolled_by=$8 WHERE udid = $1;

-- name: NewDeviceReplacingExistingResetCache :exec
DELETE FROM device_cache WHERE device_id=$1;

-- name: NewDeviceReplacingExistingResetInventory :exec
DELETE FROM device_cache WHERE device_id=$1;

-- name: SetDeviceState :exec
UPDATE devices SET state=$2 WHERE id = $1;

-- name: DeviceUserUnenrollment :exec
UPDATE devices SET state='user_unenrolled', enrollment_type='Unenrolled', azure_did='', nodecache_version='', lastseen=to_timestamp(CAST(0 as bigint)/1000), lastseen_status=0, enrolled_at=to_timestamp(CAST(0 as bigint)/1000), enrolled_by=NULL WHERE id = $1;

-- name: GetDevices :many
-- Exposed via API
SELECT id, name, model FROM devices LIMIT 100;

-- name: GetBasicDevice :one
-- Exposed via API
SELECT id, name, description, model FROM devices WHERE id = $1 LIMIT 1;

-- name: GetBasicDeviceScopedGroups :many
-- Exposed via API
SELECT groups.id, groups.name FROM groups INNER JOIN group_devices ON group_devices.group_id=groups.id WHERE group_devices.device_id = $1;

-- name: GetBasicDeviceScopedPolicies :many
-- Exposed via API
SELECT * FROM policies INNER JOIN group_policies ON group_policies.policy_id = policies.id INNER JOIN group_devices ON group_devices.group_id=group_policies.group_id WHERE group_devices.device_id = $1;

-- name: GetDevice :one
SELECT * FROM devices WHERE id = $1 LIMIT 1;

-- name: GetDeviceByUDID :one
SELECT * FROM devices WHERE udid = $1 LIMIT 1;

-- name: DeviceCheckinStatus :exec
UPDATE devices SET lastseen=NOW(), lastseen_status=$2 WHERE id = $1; -- TODO: Merge this with last checkin status

-- name: GetDevicesPayloads :many
SELECT policies_payload.* FROM group_devices INNER JOIN group_policies ON group_policies.group_id=group_devices.group_id INNER JOIN policies_payload ON policies_payload.policy_id=group_policies.policy_id WHERE group_devices.device_id = $1;

-- name: GetDevicesPayloadsAwaitingDeployment :many
SELECT id, uri, format, type, value, exec FROM group_devices INNER JOIN group_policies ON group_policies.group_id=group_devices.group_id INNER JOIN policies_payload ON policies_payload.policy_id=group_policies.policy_id WHERE group_devices.device_id = $1 AND NOT EXISTS (SELECT 1 FROM device_cache WHERE device_cache.payload_id = policies_payload.id AND device_cache.device_id=group_devices.device_id);

-- name: GetDevicesDetachedPayloads :many
SELECT id, uri, exec FROM device_cache INNER JOIN policies_payload ON policies_payload.id=device_cache.payload_id WHERE device_cache.device_id = $1 AND NOT EXISTS (SELECT policies_payload.* FROM group_devices INNER JOIN group_policies ON group_policies.group_id=group_devices.group_id INNER JOIN policies_payload ON policies_payload.policy_id=group_policies.policy_id WHERE group_devices.device_id = device_cache.device_id);

-- name: NewDeviceCacheNode :one
INSERT INTO device_cache(device_id, payload_id) VALUES ($1, $2) RETURNING cache_id;

-- name: DeleteDeviceCacheNode :exec
DELETE FROM device_cache WHERE device_id = $1 AND payload_id = $2;

-- name: UpdateDeviceInventoryNode :exec
INSERT INTO device_inventory(device_id, uri, format, value) VALUES ($1, $2, $3, $4); -- TODO: Update or Replace

-- name: GetPolicies :many
-- Exposed via API
SELECT id, name FROM policies LIMIT 100;

-- name: GetPolicy :one
-- Exposed via API
SELECT id, name, description, priority FROM policies WHERE id = $1 LIMIT 1;

-- name: GetPoliciesPayloads :many
SELECT * FROM policies_payload WHERE policy_id = $1;

-- name: Settings :one
SELECT * FROM settings LIMIT 1;

-- name: GetRawCert :one
SELECT cert, key FROM certificates WHERE id = $1 LIMIT 1;

-- name: CreateRawCert :exec
INSERT INTO certificates(id, cert, key) VALUES ($1, $2, $3);