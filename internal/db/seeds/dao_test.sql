--
-- Truncate and seed test data. Only use for testing!
--
BEGIN;

-- Truncate all tables in the integration schema
DO
$do$
  BEGIN
    EXECUTE
      (SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
       FROM pg_class
       WHERE relkind = 'r'
         AND relnamespace = 'integration'::regnamespace);
  END
$do$;

-- Seed the data
INSERT INTO accounts(id, account_number, org_id)
VALUES (1, '1', '1')
ON CONFLICT DO NOTHING;

INSERT INTO pubkeys(id, account_id, name, body, fingerprint)
VALUES (1, 1, 'lzap-ed25519-2021',
        'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap',
        'SHA256:gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk')
ON CONFLICT DO NOTHING;

-- Reset all primary key sequences. This can possibly slow down seeds in tests, in that case
-- let's use implicit primary keys.
SELECT reset_sequences('integration');

COMMIT;
