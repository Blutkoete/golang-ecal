%module ecalc
%{
/* Includes the header in the wrapper code */
#include "/usr/include/ecal/ecalc.h"
#include "/usr/include/ecal/ecal_os.h"
#include "/usr/include/ecal/ecal_defs.h"
#include "/usr/include/ecal/cimpl/ecal_callback_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_core_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_event_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_init_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_log_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_monitoring_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_process_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_publisher_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_qos_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_service_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_subscriber_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_proto_dyn_json_subscriber_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_time_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_timer_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_tlayer_cimpl.h"
#include "/usr/include/ecal/cimpl/ecal_util_cimpl.h"
%}
 
/* Parse the header file to generate wrappers */
%include "/usr/include/ecal/ecalc.h"
%include "/usr/include/ecal/ecal_os.h"
%include "/usr/include/ecal/ecal_defs.h"
%include "/usr/include/ecal/cimpl/ecal_callback_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_core_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_event_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_init_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_log_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_monitoring_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_process_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_publisher_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_qos_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_service_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_subscriber_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_proto_dyn_json_subscriber_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_time_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_timer_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_tlayer_cimpl.h"
%include "/usr/include/ecal/cimpl/ecal_util_cimpl.h"