

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("d2de45f4682384c59206042dce79a04fec7ef90e")
var objectStoreGroup = ComponentGroup ("12eb4738d08c260872f6da2980aaec4f6995f570")


type callbacks struct {
	objectStoreIp net.IP
	objectStorePort uint16
	objectStoreFqdn string
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
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
