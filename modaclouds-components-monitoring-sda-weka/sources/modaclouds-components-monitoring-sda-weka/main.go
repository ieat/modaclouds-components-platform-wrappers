

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("c79a3a7cb3c78724cc88001b67e332376cd95d0f")
var ddaGroup = ComponentGroup ("2202877ee831a07c419eb9c62721e220d3251483")
var kbGroup = ComponentGroup ("8170ac9800426eb467537d37b7172e1d96f993b7")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	ddaIp net.IP
	ddaPort uint16
	ddaFqdn string
	kbIp net.IP
	kbPort uint16
	kbFqdn string
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
	
	_server.Transcript.TraceInformation ("resolving the DDA HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (ddaGroup, "modaclouds-monitoring-dda:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.ddaIp = _ip_1
		_callbacks.ddaPort = _port_1
		_callbacks.ddaFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the DDA endpoint: `%s:%d`;", _callbacks.ddaIp.String (), _callbacks.ddaPort)
	
	_server.Transcript.TraceInformation ("resolving the knowledgebase HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (kbGroup, "modaclouds-knowledgebase:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.kbIp = _ip_1
		_callbacks.kbPort = _port_1
		_callbacks.kbFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the KB endpoint: `%s:%d`;", _callbacks.kbIp.String (), _callbacks.kbPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_monitoring_sda_weka_run")
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_MONITORING_SDA_WEKA_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_MONITORING_SDA_WEKA_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_MONITORING_DDA_ENDPOINT_IP" : _callbacks.ddaIp.String (),
			"MODACLOUDS_MONITORING_DDA_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.ddaPort),
			"MODACLOUDS_KNOWLEDGEBASE_ENDPOINT_IP" : _callbacks.kbIp.String (),
			"MODACLOUDS_KNOWLEDGEBASE_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.kbPort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-monitoring-sda-weka:get-http-endpoint" :
			
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
