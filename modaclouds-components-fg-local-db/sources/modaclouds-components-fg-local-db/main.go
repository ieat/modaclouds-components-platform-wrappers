

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("3cf6a77225877c935f3208a8d3e5eb8f455cc96b")
var objectStoreGroup = ComponentGroup ("12eb4738d08c260872f6da2980aaec4f6995f570")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	objectStoreIp net.IP
	objectStorePort uint16
	objectStoreFqdn string
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
			"MODACLOUDS_FUSEKI_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_FUSEKI_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
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
		
		case "modaclouds-fg-local-db:get-http-endpoint" :
			
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
