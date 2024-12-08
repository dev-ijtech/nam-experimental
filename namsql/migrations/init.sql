/* it's pretty self-explanatory, run this script on an empty DB to initialise it to be used by nam */

PRAGMA journal_mode = WAL;

PRAGMA busy_timeout = 5000;

PRAGMA foreign_keys = ON;

CREATE TABLE "Device" (
  "id"             INTEGER,
  "name"           TEXT NOT NULL,
  "managementIPv4" TEXT NOT NULL,
  "vendor"         TEXT NOT NULL,
  "version"        TEXT NOT NULL,
  "createdAt"      TEXT NOT NULL,
  "updatedAt"      TEXT NOT NULL,
  PRIMARY KEY ("id")
) STRICT;
