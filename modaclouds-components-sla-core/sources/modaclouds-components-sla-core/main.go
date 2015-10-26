

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("3d349decd4ea24f20c5aae4fbaa9b5024c28ab03")
var managerGroup = ComponentGroup ("70e89545c5078bb95618f0fc5ff9283c87d8e687")
var mysqlGroup = ComponentGroup ("be149e7b52c7cbe0695e208081ffaefbbc5778a7")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	managerIp net.IP
	managerPort uint16
	managerFqdn string
	mysqlIp net.IP
	mysqlPort uint16
	mysqlFqdn string
	mysqlAccount string
	mysqlPassword string
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
	_server.Transcript.TraceInformation ("acquiring the HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("http")); _error != nil {
		return _error
	} else {
		_callbacks.httpIp = _ip_1
		_callbacks.httpPort = _port_1
		_callbacks.httpFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the HTTP endpoint: `%s:%d`;", _callbacks.httpIp.String (), _callbacks.httpPort)
	
	_server.Transcript.TraceInformation ("resolving the manager HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (managerGroup, "modaclouds-monitoring-manager:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.managerIp = _ip_1
		_callbacks.managerPort = _port_1
		_callbacks.managerFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the manager endpoint: `%s:%d`;", _callbacks.managerIp.String (), _callbacks.managerPort)
	
	_server.Transcript.TraceInformation ("resolving the MySQL endpoint...")
	if _outputs_2, _error := _server.ComponentCall (mysqlGroup, "mosaic-mysql:get-sql-endpoint", nil); _error != nil {
		return _error
	} else {
		
		_outputs_1 := _outputs_2.(map[string]interface{})
		
		_ip_2 := _outputs_1["ip"].(string)
		_port_1 := uint16 (_outputs_1["port"].(float64))
		_fqdn_1 := _outputs_1["fqdn"].(string)
		
		_ip_1 := net.ParseIP (_ip_2)
		if _ip_1 == nil {
			return errors.New ("invalid IP address")
		}
		
		_callbacks.mysqlIp = _ip_1
		_callbacks.mysqlPort = _port_1
		_callbacks.mysqlFqdn = _fqdn_1
		_callbacks.mysqlAccount = _outputs_1["administrator-login"].(string)
		_callbacks.mysqlPassword = _outputs_1["administrator-password"].(string)
	}
	
	_server.Transcript.TraceInformation ("  * using the MySQL endpoint: `%s:%d`;", _callbacks.mysqlIp.String (), _callbacks.mysqlPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_service_run")
	
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_SLACORE_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_SLACORE_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_TOWER4CLOUDS_MANAGER_ENDPOINT_IP" : _callbacks.managerIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_MANAGER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.managerPort),
			"MODACLOUDS_MYSQL_ENDPOINT_IP" : _callbacks.mysqlIp.String (),
			"MODACLOUDS_MYSQL_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.mysqlPort),
			"MODACLOUDS_MYSQL_USERNAME" : _callbacks.mysqlAccount,
			"MODACLOUDS_MYSQL_PASSWORD" : _callbacks.mysqlPassword,
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-sla-core:get-http-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.httpIp.String (),
					"port" : _callbacks.httpPort,
					"fqdn" : _callbacks.httpFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.httpFqdn, _callbacks.httpPort),
			}
		
		default :
			
			_error = errors.New ("invalid-operation")
	}
	
	return _outputs, _error
}


func main () () {
	PreMain (& callbacks {}, packageTranscript)
}


var packageTranscript = transcript.NewPackageTranscript (transcript.InformationLevel)
