

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("d2de45f4682384c59206042dce79a04fec7ef90e")
var controllerGroup = ComponentGroup ("ead924d8dcfb024f365e229da1df6b29f650f9f0")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	ctlIp net.IP
	ctlPort uint16
	ctlFqdn string
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
	
	_server.Transcript.TraceInformation ("resolving the controller HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (ctlGroup, "modaclouds-load-balancer-controller:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.ctlIp = _ip_1
		_callbacks.ctlPort = _port_1
		_callbacks.ctlFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the controller endpoint: `%s:%d`;", _callbacks.ctlIp.String (), _callbacks.ctlPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_load_balancer_reasoner_run")
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_LOAD_BALANCER_REASONER_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_LOAD_BALANCER_REASONER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_LOAD_BALANCER_CONTROLLER_ENDPOINT_IP" : _callbacks.ctlIp.String (),
			"MODACLOUDS_LOAD_BALANCER_CONTROLLER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.ctlPort),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-load-balancer-reasoner:get-http-endpoint" :
			
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


var packageTranscript = transcript.NewPackageTranscript ()