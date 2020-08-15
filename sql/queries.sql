-- DO NOT RUN THIS FILE. It is used along with sqlc to generate type safe Go from SQL

-- name: GetUser :one
SELECT * FROM users WHERE upn = $1 LIMIT 1;

-- name: NewAzureADUser :exec
INSERT INTO users(upn, fullname, azuread_oid) VALUES($1, $2, $3); -- TODO: Insert or Update

-- name: NewDevice :exec
INSERT INTO devices(udid, state, enrollment_type, name, serial_number, operating_system, azure_did, enrolled_by) VALUES($1, $2, $3, $4, $5, $6, $7, $8);

-- name: NewDeviceReplacingExisting :exec
UPDATE devices SET state=$2, enrollment_type=$3, name=$4, serial_number=$5, operating_system=$6, azure_did=$7, nodecache_version='', lastseen=NOW(), lastseen_status=0, enrolled_at=NOW(), enrolled_by=$8 WHERE udid = $1;

-- name: NewDeviceReplacingExistingReset :exec
DELETE FROM device_cache WHERE device_id=$1;

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

-- name: NewDeviceCacheNode :exec
INSERT INTO device_cache(device_id, payload_id) VALUES ($1, $2); -- TODO: cache_id

-- name: DeleteDeviceCacheNode :exec
DELETE FROM device_cache WHERE device_id = $1 AND payload_id = $2;

-- name: GetPolicy :one
SELECT * FROM policies WHERE id = $1 LIMIT 1;

-- name: GetPoliciesPayloads :many
SELECT * FROM policies_payload WHERE policy_id = $1;

-- name: Settings :one
SELECT * FROM settings LIMIT 1;

-- name: GetRawCert :one
SELECT cert, key FROM certificates WHERE id = $1 LIMIT 1;

-- name: CreateRawCert :exec
INSERT INTO certificates(id, cert, key) VALUES ($1, $2, $3);