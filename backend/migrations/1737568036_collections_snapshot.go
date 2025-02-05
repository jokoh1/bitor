package migrations

import (
	"encoding/json"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// Add check if _collections table exists before validation
		var exists bool
		err := db.NewQuery("SELECT 1 FROM sqlite_master WHERE type='table' AND name='_collections'").Row(&exists)
		if err != nil {
			return err
		}

		// Skip validation if collections table doesn't exist yet
		if !exists {
			return nil
		}

		jsonData := `[
			{
				"id": "27do0wbcuyfmbmx",
				"created": "2023-03-04 20:33:00.558Z",
				"updated": "2025-01-13 19:59:12.382Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "gfk1gfjh",
						"name": "first_name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "rti2wf9g",
						"name": "last_name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "e9kpfdit",
						"name": "avatar",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "hicqyzli",
						"name": "public_ssh_key",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "dcexnmix",
						"name": "group",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "jnasf41n6wi7kse",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "dnyaib8e",
						"name": "requirePasswordChange",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != '' && (@request.auth.group.permissions.manage_users = true || @request.auth.id = id)",
				"viewRule": "@request.auth.id != '' && (@request.auth.group.permissions.manage_users = true || @request.auth.id = id)",
				"createRule": "@request.auth.group.permissions.manage_users = true",
				"updateRule": "@request.auth.id != '' && (@request.auth.group.permissions.manage_users = true || @request.auth.id = id)",
				"deleteRule": "@request.auth.group.permissions.manage_users = true",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": [],
					"manageRule": null,
					"minPasswordLength": 5,
					"onlyEmailDomains": [],
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "k80ers236gkl3vt",
				"created": "2024-10-06 15:55:19.087Z",
				"updated": "2025-01-14 16:34:49.444Z",
				"name": "api_keys",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "z7aocozw",
						"name": "name",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "piy1rfmy",
						"name": "key",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "kq69h1sb",
						"name": "key_type",
						"type": "select",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"api_key",
								"webhook_url",
								"bot_token",
								"personal_token",
								"integration_token",
								"username",
								"password",
								"access_key",
								"secret_key",
								"client_id",
								"client_secret",
								"smtp_username",
								"smtp_password",
								"webhook_id",
								"webhook_token"
							]
						}
					},
					{
						"system": false,
						"id": "rwhggekz",
						"name": "provider",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "ikf0vddc",
						"name": "key_pair_id",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "q9afc1vh",
						"name": "key_pair_type",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"single",
								"username_password",
								"access_secret",
								"client_oauth",
								"webhook_pair"
							]
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.group.permissions.manage_api_keys = true",
				"viewRule": "@request.auth.group.permissions.manage_api_keys = true",
				"createRule": "@request.auth.group.permissions.manage_api_keys = true",
				"updateRule": "@request.auth.group.permissions.manage_api_keys = true",
				"deleteRule": "@request.auth.group.permissions.manage_api_keys = true",
				"options": {}
			},
			{
				"id": "zqdmvqo2mym808a",
				"created": "2024-10-06 22:56:59.754Z",
				"updated": "2025-01-21 16:00:22.454Z",
				"name": "nuclei_scans",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "kvabi6zi",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "elcbhext",
						"name": "status",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"Manual",
								"Created",
								"Started",
								"Generating",
								"Deploying",
								"Running",
								"Finished",
								"Failed",
								"Stopped"
							]
						}
					},
					{
						"system": false,
						"id": "mqfh7qxq",
						"name": "start_time",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "zla6g5pd",
						"name": "end_time",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "wavbg80k",
						"name": "results",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "xtvjp3nh",
						"name": "error_message",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "atgqlihk",
						"name": "nuclei_profile",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "p9kapb3mla6r3i3",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "z71erxdq",
						"name": "nuclei_targets",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "mpwcusgvtr5nqxe",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "wxocunvw",
						"name": "nuclei_interact",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "qpmelfawm5975p5",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "75cyr33v",
						"name": "vm_provider",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "dl63cj8h",
						"name": "client",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "2hmr3iu22ww6uih",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "ssr9adhh",
						"name": "cron",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "uqkhiicq",
						"name": "state_bucket",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "yg0j4gv6",
						"name": "scan_bucket",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "97vf0p7o",
						"name": "api_key",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "dmti1iug",
						"name": "ip_address",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "erldsmuf",
						"name": "ansible_logs",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "6jv76vce",
						"name": "scan_profile",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4ae108li7hjsxlz",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "fambskug",
						"name": "manual_targets",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "45driikd",
						"name": "cost",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "yxl3vcof",
						"name": "vm_size",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ohuc7pvl",
						"name": "destroyed",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "n84wwn4e",
						"name": "vm_start_time",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "oaf4bxlv",
						"name": "vm_stop_time",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "idyvuxuw",
						"name": "nuclei_start_time",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "zqjodod0",
						"name": "nuclei_stop_time",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "oduiyuml",
						"name": "skipped_hosts",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "oaxz1ypa",
						"name": "archived",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "vnlapq5b",
						"name": "archived_date",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_scans\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_scans\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_scans\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_scans\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"nuclei_scans\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "4k1edieeg44pr1v",
				"created": "2024-10-23 12:55:56.281Z",
				"updated": "2024-10-23 12:55:56.281Z",
				"name": "nuclei_settings",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "whvg1kc1",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "kld5k381svrgnir",
				"created": "2024-10-23 13:54:37.831Z",
				"updated": "2024-10-23 13:54:37.831Z",
				"name": "files",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "akbuh6l5",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "n0r5k2j8",
						"name": "path",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "igrwn5fk",
						"name": "type",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "qyftku5j",
						"name": "size",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "p9kapb3mla6r3i3",
				"created": "2024-10-23 14:23:32.179Z",
				"updated": "2024-12-13 17:29:58.604Z",
				"name": "nuclei_profiles",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "so1wuqgz",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "0v3oq9hn",
						"name": "profile",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "ljtvcgpj",
						"name": "raw_yaml",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_profiles\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_profiles\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_profiles\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_profiles\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"nuclei_profiles\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "qpmelfawm5975p5",
				"created": "2024-10-23 15:47:06.874Z",
				"updated": "2024-12-13 17:29:58.600Z",
				"name": "nuclei_interact",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "uemnaws5",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "wytquf1l",
						"name": "token",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "wttmojrs",
						"name": "url",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_interact\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_interact\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_interact\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_interact\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"nuclei_interact\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "mpwcusgvtr5nqxe",
				"created": "2024-10-23 17:16:28.116Z",
				"updated": "2024-12-13 17:29:58.596Z",
				"name": "nuclei_targets",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "9emhjygh",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "tmbn3hwk",
						"name": "count",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "dxsinbwe",
						"name": "targets",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "fcatifur",
						"name": "client",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "2hmr3iu22ww6uih",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_targets\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"nuclei_targets\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_targets\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"nuclei_targets\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"nuclei_targets\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "cxzqhrd7om4n8od",
				"created": "2024-10-24 16:04:02.528Z",
				"updated": "2025-01-14 18:21:20.837Z",
				"name": "providers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "7sqo1k9g",
						"name": "name",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "kk3s1ktc",
						"name": "description",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "tchgvws3",
						"name": "provider_type",
						"type": "select",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"email",
								"slack",
								"teams",
								"discord",
								"telegram",
								"jira",
								"aws",
								"digitalocean",
								"s3"
							]
						}
					},
					{
						"system": false,
						"id": "wjmq2fay",
						"name": "key",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "k80ers236gkl3vt",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "wuogycyi",
						"name": "use",
						"type": "select",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 3,
							"values": [
								"dns",
								"compute",
								"terraform_storage",
								"scan_storage",
								"notification"
							]
						}
					},
					{
						"system": false,
						"id": "h2qal4u0",
						"name": "enabled",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "lpjqyjb4",
						"name": "settings",
						"type": "json",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.group.permissions.manage_providers = true",
				"viewRule": "@request.auth.group.permissions.manage_providers = true",
				"createRule": "@request.auth.group.permissions.manage_providers = true",
				"updateRule": "@request.auth.group.permissions.manage_providers = true",
				"deleteRule": "@request.auth.group.permissions.manage_providers = true",
				"options": {}
			},
			{
				"id": "sgc6cuzt2qx3tmo",
				"created": "2024-10-25 14:06:52.230Z",
				"updated": "2025-01-16 18:35:48.821Z",
				"name": "nuclei_results",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "z8dmj7un",
						"name": "template_id",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "jnwylirf",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "t5d5mjcy",
						"name": "description",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "7qi8hqw6",
						"name": "severity",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "nr2bnhen",
						"name": "matched_at",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "bdaf853x",
						"name": "type",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "b6v0tn7r",
						"name": "host",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ik5swohy",
						"name": "ip",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "tqzp0bxy",
						"name": "timestamp",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "6x2gjtni",
						"name": "reference",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "4qa2iqmp",
						"name": "extra_info",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "gwgqoyen",
						"name": "info",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "oozll6ix",
						"name": "port",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "katohoil",
						"name": "url",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "o8ozyqgx",
						"name": "request",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "tqipu4ew",
						"name": "response",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "4rn6efjy",
						"name": "matcher_status",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "kekejcqn",
						"name": "tags",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "mxtorbex",
						"name": "markdown",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "r3gquyif",
						"name": "client",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "2hmr3iu22ww6uih",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "tdmeffnh",
						"name": "scan_id",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "zqdmvqo2mym808a",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "cajbgmrm",
						"name": "severity_order",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "vrmjvxyv",
						"name": "acknowledged",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "s37gysv5",
						"name": "false_positive",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "136ljafk",
						"name": "curl_command",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "0k1e11zz",
						"name": "notes",
						"type": "editor",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"convertUrls": false
						}
					},
					{
						"system": false,
						"id": "lrzzhjrc",
						"name": "remediated",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "xlgezaiy",
						"name": "last_seen",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != '' && (@request.auth.group.permissions.read ?~ 'findings' || @request.auth.group.permissions.read ?~ '*')",
				"viewRule": "@request.auth.id != '' && (@request.auth.group.permissions.read ?~ 'findings' || @request.auth.group.permissions.read ?~ '*')",
				"createRule": "@request.auth.id != '' && (@request.auth.group.permissions.write ?~ 'findings' || @request.auth.group.permissions.write ?~ '*')",
				"updateRule": "@request.auth.id != '' && (@request.auth.group.permissions.write ?~ 'findings' || @request.auth.group.permissions.write ?~ '*')",
				"deleteRule": "@request.auth.id != '' && (@request.auth.group.permissions.delete ?~ 'findings' || @request.auth.group.permissions.delete ?~ '*')",
				"options": {}
			},
			{
				"id": "2hmr3iu22ww6uih",
				"created": "2024-10-25 17:58:03.331Z",
				"updated": "2025-01-08 22:42:36.500Z",
				"name": "clients",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "k8xxwlmt",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "n4eetrqg",
						"name": "hidden_name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "gtcqczd9",
						"name": "group",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "3jr0ramjgkl6uoi",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "homepage",
						"name": "homepage",
						"type": "url",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": [],
							"onlyDomains": []
						}
					},
					{
						"system": false,
						"id": "favicon",
						"name": "favicon",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [
								"image/x-icon",
								"image/vnd.microsoft.icon",
								"image/png",
								"image/jpeg",
								"image/gif"
							],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"clients\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"clients\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"clients\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"clients\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"clients\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "rx6ocb40rrc20vh",
				"created": "2024-10-30 19:57:17.115Z",
				"updated": "2024-10-30 20:14:08.698Z",
				"name": "ansible",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "qbbxu9vv",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "nhwjn6z7",
						"name": "ssh_public_key",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "4b06kszj",
						"name": "ssh_private_key",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "mu5n2c0lsdvpvip",
				"created": "2024-11-03 21:13:26.888Z",
				"updated": "2025-01-07 22:14:48.532Z",
				"name": "notifications",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "zfj4vpjg",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != ''",
				"viewRule": "@request.auth.id != ''",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			},
			{
				"id": "v924wfw2jnv7tye",
				"created": "2024-11-05 17:37:52.429Z",
				"updated": "2024-11-05 18:36:23.913Z",
				"name": "nuclei_scan_archives",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "qb2yobjy",
						"name": "client_id",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "2hmr3iu22ww6uih",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "zqloxss0",
						"name": "scan_id",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "oxaaz7ci",
						"name": "s3_full_path",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "zdrodeox",
						"name": "s3_small_path",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "gsvor4ww",
						"name": "s3_provider_id",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "3jr0ramjgkl6uoi",
				"created": "2024-11-07 19:22:23.604Z",
				"updated": "2024-12-13 17:29:58.592Z",
				"name": "client_groups",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "i0cjden0",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"client_groups\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"client_groups\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"client_groups\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"client_groups\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"client_groups\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "4ae108li7hjsxlz",
				"created": "2024-10-23 17:16:28.116Z",
				"updated": "2024-12-13 17:29:58.596Z",
				"name": "scan_profiles",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "ycqibpre",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "iwzpqfrb",
						"name": "nuclei_interact",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "qpmelfawm5975p5",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "dlp1mbtr",
						"name": "vm_provider",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "v7pasadd",
						"name": "state_bucket",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "6dvoesmk",
						"name": "scan_bucket",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "cxzqhrd7om4n8od",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "xlizvorj",
						"name": "default",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "otjzu9ig",
						"name": "vm_size",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "nuclei_targets_field",
						"name": "nuclei_targets",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "mpwcusgvtr5nqxe",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"scan_profiles\" || @request.auth.group.permissions.read ~ \"*\")",
				"viewRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.read ~ \"scan_profiles\" || @request.auth.group.permissions.read ~ \"*\")",
				"createRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"scan_profiles\" || @request.auth.group.permissions.write ~ \"*\")",
				"updateRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.write ~ \"scan_profiles\" || @request.auth.group.permissions.write ~ \"*\")",
				"deleteRule": "@request.auth.id != \"\" && (@request.auth.group.name = \"admin\" || @request.auth.group.permissions.delete ~ \"scan_profiles\" || @request.auth.group.permissions.delete ~ \"*\")",
				"options": {}
			},
			{
				"id": "64vh0u4fmvdmw8v",
				"created": "2024-11-15 18:11:07.530Z",
				"updated": "2024-12-14 15:53:44.548Z",
				"name": "scheduled_scans",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "9lti8bfj",
						"name": "scan_id",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "zqdmvqo2mym808a",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "lyj3z6qo",
						"name": "frequency",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "xvwxub7r",
						"name": "cron_expression",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "2kc0hezk",
						"name": "start_date",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "wgaut6wj",
						"name": "end_date",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "hxdqecjs",
						"name": "schedule_details",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "xtuprs0au45l8d3",
				"created": "2024-12-10 18:42:15.540Z",
				"updated": "2024-12-10 20:12:09.420Z",
				"name": "settings",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "kwtzte9l",
						"name": "setup_completed",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "5a9so80mi0eefyz",
				"created": "2024-12-13 02:12:07.709Z",
				"updated": "2024-12-13 20:50:20.138Z",
				"name": "invitations",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "cyfgg9ku",
						"name": "email",
						"type": "email",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "7flyrgwt",
						"name": "token",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "no41m0qj",
						"name": "expires",
						"type": "date",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "yjpjiwc5",
						"name": "used",
						"type": "bool",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "qdo7vfmi",
						"name": "group",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "jnasf41n6wi7kse",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "jnasf41n6wi7kse",
				"created": "2024-12-13 15:52:01.810Z",
				"updated": "2025-01-07 22:14:48.531Z",
				"name": "groups",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "4njqv6zc",
						"name": "name",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "cxxcqjpo",
						"name": "description",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "zmwndsj7",
						"name": "permissions",
						"type": "json",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != ''",
				"viewRule": "@request.auth.id != ''",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			},
			{
				"id": "bzhstcqk4q9fgy5",
				"created": "2025-01-08 22:42:36.503Z",
				"updated": "2025-01-09 01:54:26.000Z",
				"name": "notification_groups",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "xkgpdjmb",
						"name": "name",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "yuoqf6qk",
						"name": "description",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "uj1zh7mk",
						"name": "members",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "27do0wbcuyfmbmx",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "",
				"viewRule": "",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			},
			{
				"id": "notification_settings",
				"created": "2025-01-10 16:03:25.109Z",
				"updated": "2025-01-22 01:23:30.760Z",
				"name": "notification_settings",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "rules",
						"name": "rules",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 0
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.manage_notifications = true",
				"viewRule": "@request.auth.manage_notifications = true",
				"createRule": "@request.auth.manage_notifications = true",
				"updateRule": "@request.auth.manage_notifications = true",
				"deleteRule": "@request.auth.manage_notifications = true",
				"options": {}
			},
			{
				"id": "yw3iuuu4rx2d4ku",
				"created": "2025-01-14 22:53:01.278Z",
				"updated": "2025-01-18 03:00:55.251Z",
				"name": "system_settings",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "ziqxxtvr",
						"name": "scan_concurrency",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 100,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "5hzhty0s",
						"name": "auto_scan_enabled",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "qwulphvg",
						"name": "auto_scan_interval",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 168,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "wijydjzt",
						"name": "retention_period",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 365,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "eeioylvn",
						"name": "debug_mode",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "kgqr1jiw",
						"name": "encyption_test",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "dr8hsrll",
						"name": "stale_threshold_days",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "wv7o8srj",
						"name": "max_cost_per_month",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 0,
							"max": 5000,
							"noDecimal": false
						}
					}
				],
				"indexes": [],
				"listRule": "",
				"viewRule": "",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		// First import collections to create tables
		if err := daos.New(db).ImportCollections(collections, true, nil); err != nil {
			return err
		}

		// Now create initial system settings record
		_, err = db.NewQuery("INSERT INTO system_settings (id, created, updated, scan_concurrency, auto_scan_enabled, auto_scan_interval, retention_period, debug_mode) VALUES ({:id}, {:created}, {:updated}, {:scan_concurrency}, {:auto_scan_enabled}, {:auto_scan_interval}, {:retention_period}, {:debug_mode})").
			Bind(map[string]interface{}{
				"id":                 "default",
				"created":            time.Now().UTC().Format(time.RFC3339),
				"updated":            time.Now().UTC().Format(time.RFC3339),
				"scan_concurrency":   5,
				"auto_scan_enabled":  false,
				"auto_scan_interval": 24,
				"retention_period":   30,
				"debug_mode":         false,
			}).Execute()

		return err
	}, func(db dbx.Builder) error {
		return nil
	})
}
