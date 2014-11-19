

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("7e079e121d36b8cec279b3116b2ec3d9a4e36045")


type callbacks struct {
	dashboardIp net.IP
	dashboardPort uint16
	dashboardFqdn string
	queryIp net.IP
	queryPort uint16
	queryFqdn string
	pickleReceiverIp net.IP
	pickleReceiverPort uint16
	pickleReceiverFqdn string
	lineReceiverIp net.IP
	lineReceiverPort uint16
	lineReceiverFqdn string
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
	_server.Transcript.TraceInformation ("acquiring the dashboard HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("dashboard")); _error != nil {
		return _error
	} else {
		_callbacks.dashboardIp = _ip_1
		_callbacks.dashboardPort = _port_1
		_callbacks.dashboardFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the dashboard HTTP endpoint: `%s:%d`;", _callbacks.dashboardIp.String (), _callbacks.dashboardPort)
	
	_server.Transcript.TraceInformation ("acquiring the query TCP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("query")); _error != nil {
		return _error
	} else {
		_callbacks.queryIp = _ip_1
		_callbacks.queryPort = _port_1
		_callbacks.queryFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the query TCP endpoint: `%s:%d`;", _callbacks.queryIp.String (), _callbacks.queryPort)
	
	_server.Transcript.TraceInformation ("acquiring the pickle-receiver TCP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("pickle-receiver")); _error != nil {
		return _error
	} else {
		_callbacks.pickleReceiverIp = _ip_1
		_callbacks.pickleReceiverPort = _port_1
		_callbacks.pickleReceiverFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the pickle-receiver TCP endpoint: `%s:%d`;", _callbacks.pickleReceiverIp.String (), _callbacks.pickleReceiverPort)
	
	_server.Transcript.TraceInformation ("acquiring the line-receiver TCP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("line-receiver")); _error != nil {
		return _error
	} else {
		_callbacks.lineReceiverIp = _ip_1
		_callbacks.lineReceiverPort = _port_1
		_callbacks.lineReceiverFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the line-receiver TCP endpoint: `%s:%d`;", _callbacks.lineReceiverIp.String (), _callbacks.lineReceiverPort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_metric_explorer_run")
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_METRIC_EXPLORER_DASHBOARD_ENDPOINT_IP" : _callbacks.dashboardIp.String (),
			"MODACLOUDS_METRIC_EXPLORER_DASHBOARD_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.dashboardPort),
			"MODACLOUDS_METRIC_EXPLORER_QUERY_ENDPOINT_IP" : _callbacks.queryIp.String (),
			"MODACLOUDS_METRIC_EXPLORER_QUERY_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.queryPort),
			"MODACLOUDS_METRIC_EXPLORER_PICKLE_RECEIVER_ENDPOINT_IP" : _callbacks.pickleReceiverIp.String (),
			"MODACLOUDS_METRIC_EXPLORER_PICKLE_RECEIVER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.pickleReceiverPort),
			"MODACLOUDS_METRIC_EXPLORER_LINE_RECEIVER_ENDPOINT_IP" : _callbacks.lineReceiverIp.String (),
			"MODACLOUDS_METRIC_EXPLORER_LINE_RECEIVER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.lineReceiverPort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-metric-explorer:get-dashboard-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.dashboardIp.String (),
					"port" : _callbacks.dashboardPort,
					"fqdn" : _callbacks.dashboardFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.dashboardFqdn, _callbacks.dashboardPort),
			}
			
		case "modaclouds-metric-explorer:get-query-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.queryIp.String (),
					"port" : _callbacks.queryPort,
					"fqdn" : _callbacks.queryFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.queryFqdn, _callbacks.queryPort),
			}
			
		case "modaclouds-metric-explorer:get-pickle-receiver-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.pickleReceiverIp.String (),
					"port" : _callbacks.pickleReceiverPort,
					"fqdn" : _callbacks.pickleReceiverFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.pickleReceiverFqdn, _callbacks.pickleReceiverPort),
			}
			
		case "modaclouds-metric-explorer:get-line-receiver-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.lineReceiverIp.String (),
					"port" : _callbacks.lineReceiverPort,
					"fqdn" : _callbacks.lineReceiverFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.lineReceiverFqdn, _callbacks.lineReceiverPort),
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
