

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("70e89545c5078bb95618f0fc5ff9283c87d8e687")
var ddaGroup = ComponentGroup ("2202877ee831a07c419eb9c62721e220d3251483")
var historyDbGroup = ComponentGroup ("3fd8108ea06e07ae7adc86cc52e8c9560c65b3c1")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	ddaIp net.IP
	ddaPort uint16
	ddaFqdn string
	historyDbIp net.IP
	historyDbPort uint16
	historyDbFqdn string
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
	
	_server.Transcript.TraceInformation ("resolving the History DB HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (historyDbGroup, "modaclouds-monitoring-history-db:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.historyDbIp = _ip_1
		_callbacks.historyDbPort = _port_1
		_callbacks.historyDbFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the Hisory DB endpoint: `%s:%d`;", _callbacks.historyDbIp.String (), _callbacks.historyDbPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_service_run")
	
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_TOWER4CLOUDS_MANAGER_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_MANAGER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_TOWER4CLOUDS_DATA_ANALYZER_ENDPOINT_IP" : _callbacks.ddaIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_DATA_ANALYZER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.ddaPort),
			// FIXME:  Resolve and use the actual public IP!
			"MODACLOUDS_TOWER4CLOUDS_DATA_ANALYZER_ENDPOINT_IP_PUBLIC" : _callbacks.ddaIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_DATA_ANALYZER_ENDPOINT_PORT_PUBLIC" : fmt.Sprintf ("%d", _callbacks.ddaPort),
			"MODACLOUDS_TOWER4CLOUDS_RDF_HISTORY_DB_ENDPOINT_IP" : _callbacks.historyDbIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_RDF_HISTORY_DB_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.historyDbPort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-monitoring-manager:get-http-endpoint" :
			
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
