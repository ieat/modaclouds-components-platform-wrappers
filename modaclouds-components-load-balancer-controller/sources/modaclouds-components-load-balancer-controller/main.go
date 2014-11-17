

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("ead924d8dcfb024f365e229da1df6b29f650f9f0")


type callbacks struct {
	controllerIp net.IP
	controllerPort uint16
	controllerFqdn string
	gatewayIp net.IP
	gatewayPort uint16
	gatewayFqdn string
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
	_server.Transcript.TraceInformation ("acquiring the controller HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("controller")); _error != nil {
		return _error
	} else {
		_callbacks.controllerIp = _ip_1
		_callbacks.controllerPort = _port_1
		_callbacks.controllerFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the controller HTTP endpoint: `%s:%d`;", _callbacks.controllerIp.String (), _callbacks.controllerPort)
	
	_server.Transcript.TraceInformation ("acquiring the gateway TCP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("gateway")); _error != nil {
		return _error
	} else {
		_callbacks.gatewayIp = _ip_1
		_callbacks.gatewayPort = _port_1
		_callbacks.gatewayFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the gateway TCP endpoint: `%s:%d`;", _callbacks.gatewayIp.String (), _callbacks.gatewayPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_load_balancer_controller_run")
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_LOAD_BALANCER_CONTROLLER_ENDPOINT_IP" : _callbacks.controllerIp.String (),
			"MODACLOUDS_LOAD_BALANCER_CONTROLLER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.controllerPort),
			"MODACLOUDS_LOAD_BALANCER_GATEWAY_ENDPOINT_IP" : _callbacks.gatewayIp.String (),
			"MODACLOUDS_LOAD_BALANCER_GATEWAY_ENDPOINT_PORT_MIN" : fmt.Sprintf ("%d", _callbacks.gatewayPort),
			"MODACLOUDS_LOAD_BALANCER_GATEWAY_ENDPOINT_PORT_MAX" : fmt.Sprintf ("%d", _callbacks.gatewayPort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-load-balancer-controller:get-controller-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.controllerIp.String (),
					"port" : _callbacks.controllerPort,
					"fqdn" : _callbacks.controllerFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.controllerFqdn, _callbacks.controllerPort),
			}
			
		case "modaclouds-load-balancer-controller:get-gateway-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.gatewayIp.String (),
					"port" : _callbacks.gatewayPort,
					"fqdn" : _callbacks.gatewayFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.gatewayFqdn, _callbacks.gatewayPort),
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
