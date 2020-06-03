package rpc

import (
	"net/http"
	//utils "github.com/Varunram/essentials/utils"
	"io"
	"log"
	"strings"

	erpc "github.com/Varunram/essentials/rpc"
)

// ParticleRPC contains a list of all particle related endpoints
var ParticleRPC = map[int][]string{
	1:  {"/particle/devices", "GET", "accessToken"},                     // GET
	2:  {"/particle/productinfo", "GET", "accessToken", "productInfo"},  // GET
	3:  {"/particle/deviceinfo", "GET", "accessToken", "deviceId"},      // GET
	4:  {"/particle/deviceping", "GET", "accessToken", "deviceId"},      // GET
	5:  {"/particle/devicesignal", "GET", "signal", "accessToken"},      // GET
	6:  {"/particle/getdeviceid", "GET", "serialNumber", "accessToken"}, // GET
	7:  {"/particle/diag/last", "GET", "accessToken", "deviceId"},       // GET
	8:  {"/particle/diag/all", "GET", "accessToken", "deviceId"},        // GET
	9:  {"/particle/user/info", "GET", "accessToken"},                   // GET
	10: {"/particle/sims", "GET", "accessToken"},                        // GET
}

// setupParticleHandlers sets up all the particle related endpoints
func setupParticleHandlers() {
	listAllDevices()
	listProductInfo()
	getDeviceInfo()
	pingDevice()
	signalDevice()
	serialNumberInfo()
	getDiagnosticsLast()
	getAllDiagnostics()
	getParticleUserInfo()
	getAllSims()
}

// ParticleDevice is a structure to parse the returned particle.io data
type ParticleDevice struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	LastApp               string `json:"last_app"`
	LastIPAddress         string `json:"last_ip_address"`
	ProductID             int    `json:"product_id"`
	Connected             bool   `json:"connected"`
	PlatformID            int    `json:"platform_id"`
	Cellular              bool   `json:"cellular"`
	Notes                 string `json:"notes"`
	Status                string `json:"status"`
	SerialNumber          string `json:"serial_number"`
	CurrentBuildTarget    string `json:"current_build_target"`
	SystemFirmwareVersion string `json:"system_firmware_version"`
	DefaultBuildTarget    string `json:"default_build_target"`
}

// ParticleProductDevice is a structure to parse returned particle.io data
type ParticleProductDevice struct {
	ID                             string   `json:"id"`
	ProductID                      int      `json:"product_id"`
	LastIPAddress                  string   `json:"last_ip_address"`
	LastHandshakeAt                string   `json:"last_handshake_at"`
	UserID                         string   `json:"user_id"`
	Online                         bool     `json:"online"`
	Name                           string   `json:"name"`
	PlatformID                     int      `json:"platform_id"`
	FirmwareProductID              int      `json:"firmware_product_id"`
	Quarantined                    bool     `json:"quarantined"`
	Denied                         bool     `json:"denied"`
	Development                    bool     `json:"development"`
	Groups                         []string `json:"groups"`
	TargetedFirmwareReleaseVersion string   `json:"targeted_firmware_release_version"`
	SystemFirmwareVersion          string   `json:"system_firmware_version"`
	SerialNumber                   string   `json:"serial_number"`
	Owner                          string   `json:"owner"`
}

// ParticleProductInfo is a structure to parse returned particle.io data
type ParticleProductInfo struct {
	Devices []ParticleProductDevice
}

// ParticlePingResponse is a structure to parse returned particle.io data
type ParticlePingResponse struct {
	Online bool `json:"online"`
	Ok     bool `json:"ok"`
}

// SignalResponse is a structure to parse returned particle.io data
type SignalResponse struct {
	ID        string `json:"id"`
	Connected bool   `json:"connected"`
	Signaling bool   `json:"signaling"`
}

// SerialNumberResponse is a structure to parse returned particle.io data
type SerialNumberResponse struct {
	Ok         bool   `json:"ok"`
	DeviceID   string `json:"device_id"`
	PlatformID int    `json:"platform_id"`
}

// ParticleUser is a structure to parse returned particle.io data
type ParticleUser struct {
	Username        string   `json:"username"`
	SubscriptionIds []string `json:"subscription_ids"`
	AccountInfo     struct {
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		CompanyName     string `json:"company_name"`
		BusinessAccount bool   `json:"business_account"`
	} `json:"account_info"`
	TeamInvites         []string `json:"team_invites"`
	WifiDeviceCount     int      `json:"wifi_device_count"`
	CellularDeviceCount int      `json:"cellular_device_count"`
}

// ParticleEventStream is a structure to parse returned particle.io data
type ParticleEventStream struct {
	Data        string `json:"data"`
	TTL         string `json:"ttl"`
	PublishedAt string `json:"published_at"`
	Coreid      string `json:"coreid"`
}

// listAllDevices lists all the devices registered to the user holding the access token
func listAllDevices() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc(ParticleRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[1][2:], ParticleRPC[1][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		body := "https://api.particle.io/v1/devices?access_token=" + accessToken
		var x []ParticleDevice
		erpc.GetAndSendJson(w, body, x)
	})
}

// listProductInfo liusts all the producsts belonging to the user with the access token
func listProductInfo() {
	http.HandleFunc(ParticleRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[2][2:], ParticleRPC[2][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		productInfo := r.URL.Query()["productInfo"][0]

		body := "https://api.particle.io/v1/products/" + productInfo + "/devices?access_token=" + accessToken
		var x ParticleProductInfo
		erpc.GetAndSendJson(w, body, x)
	})
}

// getDeviceInfo returns the information of a device. REquires device id and the accesstoken
func getDeviceInfo() {
	http.HandleFunc(ParticleRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[3][2:], ParticleRPC[3][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceID := r.URL.Query()["deviceId"][0]

		body := "https://api.particle.io/v1/devices/" + deviceID + "?access_token=" + accessToken
		var x ParticleDevice
		erpc.GetAndSendJson(w, body, x)
	})
}

// pingDevice pings a device and sees whether its up. Could be useful to create a monitoring
// dashboard of sorts where people can see if their devices are online or not
func pingDevice() {
	http.HandleFunc(ParticleRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[4][2:], ParticleRPC[4][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceID := r.URL.Query()["deviceId"][0]
		body := "https://api.particle.io/v1/devices/" + deviceID + "/ping"
		payload := strings.NewReader("access_token=" + accessToken)

		erpc.PutAndSend(w, body, payload)
	})
}

// signalDevice sends a rainbow signal to the device and the device flashes in rainbow colors
// on receiving this signal. Can be set to on or off depending on whether we want the device to flash
// in rainbow colors or not
func signalDevice() {
	http.HandleFunc(ParticleRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[5][2:], ParticleRPC[5][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceID := r.URL.Query()["deviceId"][0]
		signal := r.URL.Query()["signal"][0]
		if signal != "on" && signal != "off" {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		var body string
		var payload io.Reader
		body = "https://api.particle.io/v1/devices/" + deviceID
		if signal == "ok" {
			payload = strings.NewReader("signal=" + "1" + "&access_token=" + accessToken)
			body += "?signal=" + "1" + "&accessToken=" + accessToken
		} else {
			payload = strings.NewReader("signal=" + "0" + "&access_token=" + accessToken)
		}

		erpc.PutAndSend(w, body, payload)
	})
}

// serialNumberInfo gets the device id of a device on recipt of the serial number
func serialNumberInfo() {
	http.HandleFunc(ParticleRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[6][2:], ParticleRPC[6][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		serialNumber := r.URL.Query()["serialNumber"][0]
		accessToken := r.URL.Query()["accessToken"][0]

		body := "https://api.particle.io/v1/serial_numbers/" + serialNumber + "?access_token=" + accessToken
		var x SerialNumberResponse
		erpc.GetAndSendJson(w, body, x)
	})
}

// getDiagnosticsLast gets a list of the last diagnostic report that belongs to the device
func getDiagnosticsLast() {
	http.HandleFunc(ParticleRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[7][2:], ParticleRPC[7][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceID := r.URL.Query()["deviceId"][0]

		body := "https://api.particle.io/v1/diagnostics/" + deviceID + "/last?access_token=" + accessToken
		erpc.GetAndSendByte(w, body)
	})
}

// getAllDiagnostics gets all the past diagnostic reports of the associated device id. Requires
// accessToken for authentication
func getAllDiagnostics() {
	http.HandleFunc(ParticleRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[8][2:], ParticleRPC[8][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceID := r.URL.Query()["deviceId"][0]

		body := "https://api.particle.io/v1/diagnostics/" + deviceID + "?access_token=" + accessToken
		erpc.GetAndSendByte(w, body)
	})
}

// getParticleUserInfo gets the information of a particular user associated with an accessToken
func getParticleUserInfo() {
	http.HandleFunc(ParticleRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[9][2:], ParticleRPC[9][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		body := "https://api.particle.io/v1/user?access_token=" + accessToken
		var x ParticleUser
		erpc.GetAndSendJson(w, body, x)
	})
}

// getAllSims gets the informatiomn of all sim card that areassociated with the particular accessToken
func getAllSims() {
	http.HandleFunc(ParticleRPC[10][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ParticleRPC[10][2:], ParticleRPC[10][1])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]

		body := "https://api.particle.io/v1/sims?access_token=" + accessToken
		erpc.GetAndSendByte(w, body)
	})
}
