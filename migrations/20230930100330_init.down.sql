DROP TRIGGER tr_add_user_to_notes_config ON ns_users;
DROP FUNCTION add_user_to_notes_config;

DROP TRIGGER tr_check_notes_limit ON ns_notes;
DROP FUNCTION check_notes_limit;

DROP TABLE ns_notes_config;
DROP TABLE ns_notes;
DROP TABLE ns_users;
DROP TABLE ns_roles;