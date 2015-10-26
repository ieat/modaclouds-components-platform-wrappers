

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("60932ee0f0deac557e2e4bf2d474c06ce669da29")
var managerGroup = ComponentGroup ("70e89545c5078bb95618f0fc5ff9283c87d8e687")
var lbCtlGroup = ComponentGroup ("ead924d8dcfb024f365e229da1df6b29f650f9f0")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	managerIp net.IP
	managerPort uint16
	managerFqdn string
	lbCtlIp net.IP
	lbCtlPort uint16
	lbCtlFqdn string
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
	
	_server.Transcript.TraceInformation ("resolving the LB controller HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (lbCtlGroup, "modaclouds-load-balancer-controller:get-controller-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.lbCtlIp = _ip_1
		_callbacks.lbCtlPort = _port_1
		_callbacks.lbCtlFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the LB controller HTTP endpoint: `%s:%d`;", _callbacks.lbCtlIp.String (), _callbacks.lbCtlPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_service_run")
	
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_MODELS_AT_RUNTIME_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_MODELS_AT_RUNTIME_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_TOWER4CLOUDS_MANAGER_ENDPOINT_IP" : _callbacks.managerIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_MANAGER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.managerPort),
			"MODACLOUDS_LOAD_BALANCER_CONTROLLER_ENDPOINT_IP" : _callbacks.lbCtlIp.String (),
			"MODACLOUDS_LOAD_BALANCER_CONTROLLER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.lbCtlPort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-models-at-runtime:get-http-endpoint" :
			
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
