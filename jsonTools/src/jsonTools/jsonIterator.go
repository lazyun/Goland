package goTools

import (
	"github.com/json-iterator/go"
	"fmt"
	//"reflect"
)

var TestData = `{"ver_name": "1.8.0", "sdk_result": "0", "last_update_time": 1530979231885, "msg": {"custom": "android_lucky_e1b22389c09a9298"}, "did3": "75E1DF:38EEF8:B93096", "store": "luckincoffee", "gsm.sdk.version": "8", "sdk_ver": "v5.2.2", "did1": "3C031D:1CADED:90D3D7", "sdk_times": "1", "did2": "6DF8FB:823A9B:F1BECF", "cid": "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAJm5CIFR0L4tK0LKA61YbwUcc/TSIleipCHqnjVY8zAAVDM8OKQByovvWWPLxTqb5VALnjqSk2dPWDpPTj3YNTECAwEAAQ==", "srt": 5.761262893676758, "ext_info": {"DEVICE": "cancro", "ril.subscription.types": "RUIM", "networksubstate": "LTE", "ID": "MMB29M", "pressure": "6,1002.299133,0.000000,0.000000,", "net.hostname": "android-6634e7b3c3650f4b", "gsm.version.ril-impl": "Qualcomm RIL 1.0", "battery_infos": {"level": 25, "charging": 3, "plugged": 0, "voltage": 3626, "temperature": 411, "scale": 100}, "PRODUCT": "cancro_wc_lte", "CPU_ABI2": "armeabi", "ro.ril.oem.imei1": "867831029982089", "signatures": "3082031f30820207a00302010202045d02193f300d06092a864886f70d01010b", "wifi.interface": "wlan0", "FINGERPRINT": "Xiaomi/cancro_wc_lte/cancro:6.0.1/MMB29M/8.4.28:user/release-keys", "ro.build.host": "c3-miui-ota-bd146.bj", "gsm.sim.state": "READY", "ppid": 314, "mobilestate": 13, "ro.product.manufacturer": "Huawei", "ro.product.cpu.abi": "armeabi-v7a", "ro.build.characteristics": "nosdcard", "ro.ril.miui.imei0": "954000951334156", "apkSize": 13243230, "ro.ril.oem.imei": "867831029982089", "init.svc.adbd": "stopped", "keyguard.no_require_sim": "true", "ro.secure": "1", "ro.debuggable": "0", "ril.iccid.sim1": "89860029945178976791", "persist.sys.timezone": "Asia/Shanghai", "space": {"ram": "1893492,71956", "rom": "12.338039,8.822601"}, "gsm.network.type": "GSM", "ro.build.id": "MMB29M", "persist.radio.data.iccid": "89860029945178976791", "orientation": "3,0.000000,-0.000000,0.000000,", "ro.product.cpu.abi2": "armeabi", "apkID": "2379f9ded80e6aa5780a0eaa582fd694", "ro.product.model": "HUAWEI NXT-AL10", "ro.build.date": "Fri Jun  8 22:34:05 CST 2018", "pnumber256": "FC0B7F4DFD88987F40F8308F956231EE2DEEFE2BE8075694B6610BCA8BBF4EC4", "ro.vendor.extension_library": "libqti-perfd-client.so", "mac_framework": "00:09:e8:17:31:ee", "ro.product.brand": "Huawei", "curtime": 1530979246, "gsm.operator.numeric": "46000", "persist.radio.imei1": "867831029982089", "settings_info": {"lock_pattern_autolock": "0", "screen_brightness_mode": "1", "usb_mass_storage_enabled": "1", "sound_effects_enabled": "0", "accelerometer_rotation": "0", "screen_off_timeout": "60000", "lock_pattern_visible_pattern": "0"}, "mobilestate2": 1, "BRAND": "Huawei", "wlan.driver.status": "unloaded", "ro.build.display.id": "MMB29M", "display": {"width": 1080, "height": 1920}, "date": "Sun Jul  8 00:00:46 CST 2018\n", "USER": "builder", "ro.runtime.firstboot": "1530932163885", "persist.radio.imei": "867831029982089", "screen_on": 1, "ro.product.device": "cancro", "neighboring_cellinfos": {}, "ro.hardware": "qcom", "HOST": "c3-miui-ota-bd146.bj", "CPU_ABI": "armeabi-v7a", "ro.sys.oem.sno": "6D5830A02640", "ro.build.version.sdk": "23", "pid": 3561, "ro.board.platform": "msm8974", "ro.opengles.version": "196608", "ro.carrier": "unknown", "bssid": "00:02:00:00:00:00", "ro.build.version.release": "6.0.1", "rild.libargs": "-d /dev/smd0", "ril.iccid.sim2": "89860029945178976791", "bluetooth": "00:06:5b:e5:11:bc", "ro.build.fingerprint": "Xiaomi/cancro_wc_lte/cancro:6.0.1/MMB29M/8.4.28:user/release-keys", "SERIAL": "9911543d3dbfab05", "ssid": "CMCC", "ro.product.board": "MSM8974", "TYPE": "user", "step_counter": "19,0.000000,", "ICCID": "89860029945178976791", "gsm.version.baseband": "MPSS.DI.4.0-eaa9d90", "disk": {"2": {"used": "0.0K", "fs": "/sys/fs/cgroup", "free": "909.6M", "size": "909.6M"}, "3": {"used": "0.0K", "fs": "/var", "free": "909.6M", "size": "909.6M"}, "0": {"used": "120.0K", "fs": "/dev", "free": "799.9M", "size": "800.0M"}, "10": {"used": "0.0K", "fs": "/storage", "free": "909.6M", "size": "909.6M"}, "11": {"used": "3.4G", "fs": "/storage/emulated", "free": "8.9G", "size": "12.3G"}, "8": {"used": "4.1M", "fs": "/persist", "free": "11.6M", "size": "15.7M"}, "4": {"used": "0.0K", "fs": "/mnt", "free": "909.6M", "size": "909.6M"}, "6": {"used": "3.4G", "fs": "/data", "free": "8.9G", "size": "12.3G"}, "5": {"used": "1.1G", "fs": "/system", "free": "88.1M", "size": "1.2G"}, "9": {"used": "53.1M", "fs": "/firmware", "free": "10.9M", "size": "64.0M"}, "1": {"used": "0.0K", "fs": "/sys/fs/cgroup", "free": "909.6M", "size": "909.6M"}, "12": {"used": "0.0K", "fs": "/storage/self", "free": "909.6M", "size": "909.6M"}, "7": {"used": "6.2M", "fs": "/cache", "free": "371.5M", "size": "377.8M"}}, "root": 0, "audio_mode": 1, "apkPath": "/data/app/com.lucky.luckyclient-1/base.apk", "pkglist": "com.android.onetimeinitializer,", "apk_debug": 0, "cpuinfo": {"Hardware": "Qualcomm MSM8974PRO", "processor": "0"}, "ro.build.type": "user", "TAGS": "release-keys", "networkstate": "MOBILE", "ro.ril.miui.imei1": "954000951334156", "ro.adb.secure": "1", "ro.boot.serialno": "9911543d3dbfab05", "BOOTLOADER": "unknown", "ro.build.description": "cancro_wc_lte-user 6.0.1 MMB29M 8.4.28 release-keys", "BOARD": "MSM8974", "gsm.current.phone-type": "1", "gravity": "9,-0.168007,0.307052,9.800402,", "HARDWARE": "qcom", "fonts": "CDFEFE824234BE6C60CAF2925542B051", "fme": "88178,248,31952,21,106606,209,80709,214,71768,225,", "IMEI": "954000951334156", "sd_cid": "15010041474e44335201cbfbe23c8200\n", "accelerometer": "1,-0.177063,-0.061020,10.191025,", "ro.qc.sdk.gestures.camera": "false", "ro.build.date.utc": "1528468445", "audio_volume": 0, "IMEI_PHONE": "954000951334156", "ro.qc.sdk.sensors.gestures": "true", "signatures_md5": "13B13AB99520877C8FCEA3E630E6D135", "ro.serialno": "9911543d3dbfab05", "screen_brightness": 102, "IMSI": "460002994517897", "MANUFACTURER": "Huawei", "self_id": "1fadfbda997449a08cae118adfb6367f", "MODEL": "HUAWEI NXT-AL10", "rotation_vector": "11,0.000000,0.000000,0.000000,1.000000,3.141593,", "launcher": "com.miui.home", "network": {"socket": "cnd,cryptd,dnsproxyd,dpmd,dpmwrapper,fwmarkd,ims_datad,ims_qmid,installd,lmkd,logd,logdr,logdw,mdns,mtpd,netd,netmgr,nims,perfd,pps,property_service,qmux_audio,qmux_bluetooth,qmux_gps,qmux_radio,rild,rild-debug,sap_uim_socket1,sensor_ctl_socket,thermal-recv-client,thermal-recv-passive-client,thermal-send-client,vold,zygote,"}, "DISPLAY": "MMB29M", "ro.nfc.port": "I2C", "android_id": "6634e7b3c3650f4b", "gsm.operator.isroaming": "true", "trace": 0, "id": "uid=10129(u0_a129) gid=10129(u0_a129) groups=10129(u0_a129),1007(log),3001(net_bt_admin),3002(net_bt),3003(inet),9997(everybody),50129(all_a129) context=u:r:untrusted_app:s0:c512,c768\n", "ro.com.google.clientidbase": "android-xiaomi", "multi_app_path": "/data/user/0/com.lucky.luckyclient/app_DUHOME", "ro.gps.agps_provider": "1"}, "app_label": "luckin coffee", "did": "3C031D:1CADED:90D3D7", "ldid": "DuKx0cCMaKfpLcS90GG212F8+DtkJuKljMYxrkYhAiyIDCDGdl4zElSmfwypFuXVgrV5piF3fZtKdb5PLzE2vD7g", "ver_code": 22, "pkg_name": "com.lucky.luckyclient", "lnt": "0.032372", "na_rseq": 0, "na_seq": 1, "first_install_time": 1530979231885, "mac_addr": "00:09:e8:17:31:ee"}`

func JsonIterGetOneFromByte() {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	ret := jsoniter.Get(val, "Colors", 0).ToString()
	fmt.Println("JsonIterGetOneFromByte ret is", ret)
}


func JsonIterMashal() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var jsonRet map[string]interface{}
	err := json.Unmarshal([]byte(TestData), &jsonRet)
	if nil != err {
		fmt.Println("Json Iterator marshal err",  err)
	}

	getValue := jsonRet["pkg_name"]
	fmt.Println("Json data get pkg_name value is", getValue )
}


func JsonIterGet() {
	val := []byte(TestData)
	ret := jsoniter.Get(val, "msg", "custom").ToString()
	fmt.Println("JsonIterGetOneFromByte get msg.custom ret is", ret)

	ret = jsoniter.Get(val, "msg").ToString()
	fmt.Println("JsonIterGetOneFromByte get msg ret is", ret)
}


func JsonIterLoad(jsonString string) (error, func(key_path ...interface{}) (jsoniter.Any)) {
	var jsonRet map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(TestData), &jsonRet)
	if nil != err {
		return err, nil
	}

	varByte := []byte(jsonString)
	return nil, func(key_path ...interface{}) (jsoniter.Any) {

		fmt.Println("emm", key_path)
		return jsoniter.Get(varByte, key_path ...)
	}

}