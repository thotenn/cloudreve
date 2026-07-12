package constants

// These values will be injected at build time, DO NOT EDIT.

// BackendVersion 当前后端版本号
// Bumped to 4.14.1 in this fork to trigger the schema migration that creates the
// group⇄storage_policy allowed-set join table and backfills it (see inventory/migration.go).
// Bumped to 4.14.2 (APP-101) to trigger the migration that creates the
// media_process_task table and seeds the media-compression default settings.
//
// Fork versioning track: "<upstream_next_patch>-thotenn.<n>", offset one patch above the synced
// upstream so it is BOTH a unique marker string AND semver-greater than upstream's own marker.
// Two mechanisms depend on it (inventory/migration.go): (1) needMigration() keys on the EXACT
// string db_version_<this> — a fork suffix never collides with upstream's db_version_<upstream>
// that any upstream-migrated DB already holds, so the fork's additive migration always fires;
// (2) applyPatches() parses it as semver (must be valid or migrate errors) and gates data patches
// on max(stored markers) < patch.EndVersion — being > upstream keeps future fork patches orderable.
// Bump <n> whenever the fork adds/changes schema; sk-updater resets the base to the synced
// upstream (+1 patch) on each sync and restarts <n> at 1. Code is synced to upstream 4.17.0.
var BackendVersion = "4.17.1-thotenn.1"

// IsPro 是否为Pro版本
var IsPro = "false"

var IsProBool = IsPro == "true"

// LastCommit 最后commit id
var LastCommit = "000000"

const (
	APIPrefix      = "/api/v4"
	APIPrefixSlave = "/api/v4/slave"
	CrHeaderPrefix = "X-Cr-"
)

const CloudreveScheme = "cloudreve"

type (
	FileSystemType string
)

const (
	FileSystemMy           = FileSystemType("my")
	FileSystemShare        = FileSystemType("share")
	FileSystemTrash        = FileSystemType("trash")
	FileSystemSharedWithMe = FileSystemType("shared_with_me")
	FileSystemUnknown      = FileSystemType("unknown")
)
