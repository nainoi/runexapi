package config

const DB_HOST = "mongodb://mongodb:27017"
const DB_NAME = "runex"
const db_pass = "idever@987"
const PORT_WEB_SERVICE = ":3006"
const ID_KEY = "id"
const ROLE_KEY = "role"
const PF = "pf"
const SECRET_KEY = "RUN#987%Ex@IdevEr"
const RE_SECRET_KEY = "RUN#987%Ex@IdevErThink@2019"

const UPLOAD_IMAGE = "/upload/image/event/"
const UPLOAD_AVATAR = "/upload/image/profile/"
const UPLOAD_EVENT = "/upload/image/event/"
const HTTP = "http://"
const HTTPS = "https://"
const PrefixUrl = "http://localhost:3306"
const ImageSavePath = "upload/image/slip/"
const ImageMaxSize = 5
const RuntimeRootPath = "runtime/"
const (
	ADMIN     string = "ADMIN"
	EVENTER   string = "EVENTER"
	MEMBER    string = "MEMBER"
	SUPERUSER string = "SUPERUSER"
)

const (
	PAYMENT_WAITING         string = "PAYMENT_WAITING"
	PAYMENT_WAITING_APPROVE string = "PAYMENT_WAITING_APPROVE"
	PAYMENT_SUCCESS         string = "PAYMENT_SUCCESS"
	PAYMENT_FAIL            string = "PAYMENT_FAIL"

	PAYMENT_TRANSFER       string = "PAYMENT_TRANSFER"
	PAYMENT_CREDIT_CARD    string = "PAYMENT_CREDIT_CARD"
	PAYMENT_ONLINE_BANKING string = "PAYMENT_ONLINE_BANKING"
	PAYMENT_QRCODE         string = "PAYMENT_QRCODE"
	PAYMENT_FREE           string = "PAYMENT_FREE"
)
