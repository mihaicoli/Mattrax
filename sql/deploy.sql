-- A Basic Deployment of Mattrax
-- This should be run after the `schema.sql` file. DO NOT RUN the `queries.sql` file, it is used internally.
-- This file will be removed once the UI and code tests replace its functionality

INSERT INTO settings VALUES ('Acme School Inc', 'oscar@acme.otbeaumont.me', DEFAULT, DEFAULT,  '8cfc652d-4e83-4dda-ad8e-02c1660b807e');
INSERT INTO users VALUES ('oscar@otbeaumont.me', 'Oscar Beaumont');

INSERT INTO groups VALUES ('1', 'Test Devices');

INSERT INTO policies VALUES ('1', 'Baseline');
INSERT INTO policies_payload VALUES (DEFAULT, '1', './Vendor/MSFT/Policy/Config/WindowsLogon/EnableFirstLogonAnimation', 'int', DEFAULT, '0');
INSERT INTO group_policies VALUES ('1', '1');

INSERT INTO policies VALUES ('2', 'Student Restrictions');
INSERT INTO policies_payload VALUES (DEFAULT, '2', './Vendor/MSFT/Policy/Config/Camera/AllowCamera', 'int', DEFAULT, '0');
INSERT INTO policies_payload VALUES (DEFAULT, '2', './Device/Vendor/MSFT/Policy/Config/Connectivity/AllowBluetooth', 'int', DEFAULT, '0');

INSERT INTO policies VALUES ('3', 'Local Admin Account');
INSERT INTO policies_payload VALUES (DEFAULT, '3', './Device/Vendor/MSFT/Accounts/Users/mttx/Password', 'chr', DEFAULT, 'password');

-- Primitive Application Management. Will be replaced with propper application management in future update!
INSERT INTO policies VALUES ('4', '[TEMP] Spotify MBES Application');
INSERT INTO policies_payload VALUES (DEFAULT, '4', './User/Vendor/MSFT/EnterpriseModernAppManagement/AppInstallation/SpotifyAB.SpotifyMusic_zpdnekdrzrea0/StoreInstall', 'xml', 'text/plain', '<Application id="9NCBCSZSJRSB" flags="0" skuid="0016" />', TRUE);