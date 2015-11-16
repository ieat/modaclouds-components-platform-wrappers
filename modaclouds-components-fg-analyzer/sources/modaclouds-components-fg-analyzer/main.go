

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("283dc6bea50ff0fed89b36860bad571fe3541780")
var fusekiGroup = ComponentGroup ("8170ac9800426eb467537d37b7172e1d96f993b7")
var objectStoreGroup = ComponentGroup ("12eb4738d08c260872f6da2980aaec4f6995f570")


type callbacks struct {
	fusekiIp net.IP
	fusekiPort uint16
	fusekiFqdn string
	objectStoreIp net.IP
	objectStorePort uint16
	objectStoreFqdn string
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
	_server.Transcript.TraceInformation ("resolving the Fuseki HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (fusekiGroup, "modaclouds-fuseki:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.fusekiIp = _ip_1
		_callbacks.fusekiPort = _port_1
		_callbacks.fusekiFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the Fuseki HTTP endpoint: `%s:%d`;", _callbacks.fusekiIp.String (), _callbacks.fusekiPort)
	
	_server.Transcript.TraceInformation ("resolving the object store HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (objectStoreGroup, "mosaic-object-store:get-service-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.objectStoreIp = _ip_1
		_callbacks.objectStorePort = _port_1
		_callbacks.objectStoreFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the object store HTTP endpoint: `%s:%d`;", _callbacks.objectStoreIp.String (), _callbacks.objectStorePort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_service_run")
	
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_FUSEKI_ENDPOINT_IP" : _callbacks.fusekiIp.String (),
			"MODACLOUDS_FUSEKI_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.fusekiPort),
			"MOSAIC_OBJECT_STORE_ENDPOINT_IP" : _callbacks.objectStoreIp.String (),
			"MOSAIC_OBJECT_STORE_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.objectStorePort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		default :
			
			_error = errors.New ("invalid-operation")
	}
	
	return _outputs, _error
}


func main () () {
	PreMain (& callbacks {}, packageTranscript)
}


var packageTranscript = transcript.NewPackageTranscript (transcript.InformationLevel)
