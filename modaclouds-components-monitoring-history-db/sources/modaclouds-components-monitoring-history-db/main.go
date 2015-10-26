

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("3fd8108ea06e07ae7adc86cc52e8c9560c65b3c1")
var fusekiGroup = ComponentGroup ("8170ac9800426eb467537d37b7172e1d96f993b7")
var rabbitmqGroup = ComponentGroup ("8cd74b5e4ecd322fd7bbfc762ed6cf7d601eede8")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	fusekiIp net.IP
	fusekiPort uint16
	fusekiFqdn string
	rabbitmqIp net.IP
	rabbitmqPort uint16
	rabbitmqFqdn string
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
	
	_server.Transcript.TraceInformation ("resolving the Fuseki HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (fusekiGroup, "modaclouds-fuseki:get-http-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.fusekiIp = _ip_1
		_callbacks.fusekiPort = _port_1
		_callbacks.fusekiFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the Fuseki HTTP endpoint: `%s:%d`;", _callbacks.fusekiIp.String (), _callbacks.fusekiPort)
	
	_server.Transcript.TraceInformation ("resolving the RabbitMQ AMQP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (rabbitmqGroup, "mosaic-rabbitmq:get-broker-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.rabbitmqIp = _ip_1
		_callbacks.rabbitmqPort = _port_1
		_callbacks.rabbitmqFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the RabbitMQ AMQP endpoint: `%s:%d`;", _callbacks.rabbitmqIp.String (), _callbacks.rabbitmqPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_service_run")
	
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_TOWER4CLOUDS_RDF_HISTORY_DB_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_TOWER4CLOUDS_RDF_HISTORY_DB_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_FUSEKI_ENDPOINT_IP" : _callbacks.fusekiIp.String (),
			"MODACLOUDS_FUSEKI_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.fusekiPort),
			"MODACLOUDS_RABBITMQ_ENDPOINT_IP" : _callbacks.rabbitmqIp.String (),
			"MODACLOUDS_RABBITMQ_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.rabbitmqPort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-monitoring-history-db:get-http-endpoint" :
			
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
