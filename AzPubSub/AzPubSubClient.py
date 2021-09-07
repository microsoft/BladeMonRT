'''
Copyright (c) Microsoft Corporation.  All rights reserved.

Module for publishing and consuming events from AzPubSub(wrapper around azpubsub.dll)
'''

import os
import logging
import json
from ctypes import *
from ctypes.wintypes import *
from enum import IntEnum
import win32api
from abc import abstractmethod, ABC

script_dir = os.path.abspath(os.path.dirname(__file__))
dll_path = script_dir

# updating path so that python interpreter can locate AzPubSub.dll and all its dependencies
os.environ["PATH"] += os.pathsep + dll_path

# use windll as library functions are exported as __stdcall
AzPubSub = windll.LoadLibrary("AzPubSub.dll")

'''windows types'''
NULL = None
LONG64 = c_longlong
S_OK = 0
PHANDLE = POINTER(HANDLE)
LPPSTR = POINTER(LPSTR)
ERROR_MORE_DATA = 234
VOID = None
ENUM = c_uint

'''AzPubSub.dll types'''
KAFKA_ERROR_CODE = c_int64
LOG_LEVEL = c_int64
EVENT_TYPE = c_int64
HCLIENT = HANDLE
HCONFIG = HANDLE
HMESSAGE = HANDLE
HPRODUCER = HANDLE
HPRODUCERTOPIC = HANDLE
HRESPONSE = HANDLE
PHRESPONSE = POINTER(HANDLE) # this is not same as POINTER(HRESPONSE)


SIMPLE_PRODUCER_SUCCESS_STATUS = 200
SIMPLE_PRODUCER_SUCCESS_SUB_STATUS_OK = 2000
SIMPLE_PRODUCER_SUCCESS_MESSAGE = "All messages are sent"

class AZPUBSUB_CONFIGURATION_TYPE(IntEnum):
  GLOBAL = 0
  TOPIC = 1
  SIMPLE = 2

class AZPUBSUB_SECURITY_TYPE(IntEnum):
  NONE = 0
  SSL = 1
  TOKEN = 2

class AZPUBSUB_SECURITY_FLAGS(IntEnum):
  NONE = 0
  LOCAL = 1

class AZPUBSUB_RESPONSE_OPERATION(IntEnum):
  GET_STATUS_CODE = 2000
  GET_MESSAGE = 2001
  GET_SUB_STATUS = 2002

class AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES(IntEnum):
  NONE = 0
  LOWLATENCY = 1
  MEDIUMLATENCY = 2
  HIGHLATENCY = 3

class AZPUBSUB_MESSAGE_HEADER(ctypes.Structure):
  _fields_ = [("pszKey", LPCSTR),
              ("pValue", PBYTE),
              ("cbValue", DWORD)]

'''
AzPubSub.dll function prototypes
Refer Azure\Compute-Move\src\Services\Kafka\\UpsClients\Cpp\dll\\azpubsub\AzPubSub.h
'''

AZPUBSUB_LOG_CALLBACK = CFUNCTYPE(VOID,# return type
                                  LOG_LEVEL, LPCSTR, LPVOID) # args type

AZPUBSUB_MESSAGE_CALLBACK = CFUNCTYPE(VOID,  # return type
                                      HMESSAGE, LPVOID)  # args type

AZPUBSUB_TOPIC_PARTITIONER = CFUNCTYPE(DWORD,  # return type
                                       HPRODUCERTOPIC, PBYTE, DWORD, DWORD, LPVOID)  # args type

AZPUBSUB_EVENT_CALLBACK = CFUNCTYPE(VOID,  # return type
                                    EVENT_TYPE, DWORD, LPCSTR, LOG_LEVEL, LPCSTR, LPCSTR, DWORD, LPCSTR, DWORD)  # args type

AZPUBSUB_METRICS_CALLBACK = CFUNCTYPE(VOID,  # return type
                                      LPCSTR, LONG64)  # args type

AzPubSub.AzPubSubErrorToString.argtypes = [KAFKA_ERROR_CODE, LPPSTR, PDWORD]
AzPubSub.AzPubSubErrorToString.restype = DWORD

AzPubSub.AzPubSubClientInitialize.argtypes = [AZPUBSUB_LOG_CALLBACK, LPVOID]
AzPubSub.AzPubSubClientInitialize.restype = HCLIENT

AzPubSub.AzPubSubClientClose.argtypes = [HCLIENT]
AzPubSub.AzPubSubClientClose.restype = DWORD

AzPubSub.AzPubSubCreateConfiguration.argtypes = [HCLIENT, ENUM, c_uint]
AzPubSub.AzPubSubCreateConfiguration.restype = HCONFIG

AzPubSub.AzPubSubFreeConfiguration.argtypes = [HCONFIG]
AzPubSub.AzPubSubFreeConfiguration.restype = DWORD

AzPubSub.AzPubSubSetConnection.argtypes = [HCONFIG, LPCSTR, LPCSTR, LPCSTR, LPCSTR, ENUM, DWORD]
AzPubSub.AzPubSubSetConnection.restype = DWORD

AzPubSub.AzPubSubAddStringConfiguration.argtypes = [HCONFIG, LPCSTR, LPCSTR]
AzPubSub.AzPubSubAddStringConfiguration.restype = DWORD

AzPubSub.AzPubSubAddBinaryConfiguration.argtypes = [HCONFIG, LPCSTR, c_bool]
AzPubSub.AzPubSubAddBinaryConfiguration.restype = DWORD

AzPubSub.AzPubSubAddIntegerConfiguration.argtypes = [HCONFIG, LPCSTR, c_int]
AzPubSub.AzPubSubAddIntegerConfiguration.restype = DWORD

AzPubSub.AzPubSubGetStringConfiguration.argtypes = [HCONFIG, LPCSTR, LPPSTR, DWORD, PDWORD]
AzPubSub.AzPubSubGetStringConfiguration.restype = DWORD

AzPubSub.AzPubSubOpenProducer.argtypes = [HCONFIG, AZPUBSUB_EVENT_CALLBACK, AZPUBSUB_METRICS_CALLBACK]
AzPubSub.AzPubSubOpenProducer.restype = HPRODUCER

AzPubSub.AzPubSubOpenSimpleProducer.argtypes = [HCONFIG, ENUM, LPCSTR, LPCSTR, LPCSTR]
AzPubSub.AzPubSubOpenSimpleProducer.restype = HPRODUCER

AzPubSub.AzPubSubSendMessage.argtypes = [HPRODUCER, HPRODUCERTOPIC, LPCSTR, DWORD, LPCSTR, c_size_t, c_int32, c_int64,
                                         LPVOID, AZPUBSUB_MESSAGE_CALLBACK, LPVOID]
AzPubSub.AzPubSubSendMessage.restype = DWORD

AzPubSub.AzPubSubSendMessageEx.argtypes = [HPRODUCER, LPCSTR, LPCSTR, PINT, LPCSTR, c_size_t, PHRESPONSE]
AzPubSub.AzPubSubSendMessageEx.restype = DWORD

AzPubSub.AzPubSubProducerClose.argtypes = [HPRODUCER]
AzPubSub.AzPubSubProducerClose.restype = DWORD

AzPubSub.AzPubSubMessageGetErrorMessage.argtypes = [HMESSAGE, LPPSTR, DWORD, PDWORD]
AzPubSub.AzPubSubMessageGetErrorMessage.restype = DWORD

AzPubSub.AzPubSubMessageGetErrorCode.argtypes = [HMESSAGE, PINT]
AzPubSub.AzPubSubMessageGetErrorCode.restype = DWORD

AzPubSub.AzPubSubResponseGetStatusCode.argtypes = [HRESPONSE, PINT]
AzPubSub.AzPubSubResponseGetStatusCode.restype = DWORD

AzPubSub.AzPubSubResponseGetMessage.argtypes = [HRESPONSE, LPPSTR, DWORD, PDWORD]
AzPubSub.AzPubSubResponseGetMessage.restype = DWORD

AzPubSub.AzPubSubResponseGetSubStatusCode.argtypes = [HRESPONSE, PINT]
AzPubSub.AzPubSubResponseGetSubStatusCode.restype = DWORD

AzPubSub.AzPubSubResponseClose.argtypes = [HRESPONSE]
AzPubSub.AzPubSubResponseClose.restype = DWORD

AzPubSub.AzPubSubOpenProducerTopic.argtypes = [HPRODUCER, LPCSTR, HCONFIG, AZPUBSUB_TOPIC_PARTITIONER]
AzPubSub.AzPubSubOpenProducerTopic.restype = HPRODUCERTOPIC

AzPubSub.AzPubSubCloseProducerTopic.argtypes = [HPRODUCERTOPIC]
AzPubSub.AzPubSubCloseProducerTopic.restype = DWORD

def logger_callback(level, msg, context):
  logging.info("level={0} {1}".format(level, msg.decode("ascii").strip() if msg else ""))

p_logger_callback = AZPUBSUB_LOG_CALLBACK(logger_callback)

def message_callback(h_message, context):
  message_err = c_long()
  err = AzPubSub.AzPubSubMessageGetErrorCode(h_message, byref(message_err))
  context_str = (cast(context, LPCSTR).value) if context else ""

  if err != S_OK:
    msg = "Message callback failure. AzPubSubMessageGetErrorCode failed with error: {0}".format(win32api.FormatMessage(err))
    logging.error(msg)
    send_message_callback(context_str, GlobalResponse(status_code=message_err.value, message=msg))
    return
    
  msg_len = 256
  msg = LPSTR(b'0'*msg_len)
  err = AzPubSub.AzPubSubMessageGetErrorMessage(h_message,
                                                byref(msg),
                                                msg_len,
                                                PDWORD())
  if err != S_OK:
    msg = "Message callback failure. AzPubSubMessageGetErrorMessage failed with error: {0}".format(win32api.FormatMessage(err))
    logging.error(msg)
    send_message_callback(context_str, GlobalResponse(status_code=message_err.value, message=msg))
    return

  decoded_msg = msg.value.decode("ascii")
  
  if message_err.value == S_OK:
    logging.info("Message callback success.")
    send_message_callback(context_str, GlobalResponse(status_code=message_err.value, message=decoded_msg))
    return

  logging.error("Message callback failure. Code={0} Msg:{1} Context:{2}".format(
    message_err, 
    decoded_msg, 
    context_str.decode('ascii')))

  send_message_callback(context_str, GlobalResponse(status_code=message_err.value, message=decoded_msg))


p_message_callback = AZPUBSUB_MESSAGE_CALLBACK(message_callback)

def send_message_callback(context, global_response):
  pass

class AzPubSubResponse(ABC):
  """Abstract superclass for all responses from AzPubSub
  """
  @property
  @abstractmethod
  def message(self):
    pass

  @property
  @abstractmethod
  def status_code(self):
    pass
  
  @abstractmethod
  def __bool__(self):
    pass

class SimpleResponse(AzPubSubResponse):
  """Class for parsing response handle returned from simple_producer's send_message API.
  """
  def __init__(self, hresponse):
    self.logger = logging.getLogger('Response')
    self.hresponse = hresponse
    self._status_code = self._get_response_status_code()
    self._message = self._get_response_message()
    self.sub_status_code = self._get_response_sub_status_code()

  @property
  def message(self):
    return self._message

  @property
  def status_code(self):
    return self._status_code

  def _get_response_message(self):
    buffer_length = DWORD(0)
    p_buffer_length = LPDWORD(buffer_length)

    err = AzPubSub.AzPubSubResponseGetMessage(self.hresponse,
                                     NULL,
                                     buffer_length,
                                     p_buffer_length)
    if err != ERROR_MORE_DATA:
      self.logger.error("First AzPubSubResponseGetMessage failed with: {0}".format(win32api.FormatMessage(err)))
      return None

    p_msg_buffer = LPSTR(b'0'*buffer_length.value)
    pp_msg_buffer = pointer(p_msg_buffer)

    err = AzPubSub.AzPubSubResponseGetMessage(self.hresponse,
                                     pp_msg_buffer,
                                     buffer_length,
                                     p_buffer_length)
    if err != S_OK:
      self.logger.error("Second AzPubSubResponseGetMessage failed with: {0}".format(win32api.FormatMessage(err)))
      return None

    return p_msg_buffer.value.decode("ascii")

  def _get_response_status_code(self):
    value = c_int(0)
    p_value = pointer(value)
    err = AzPubSub.AzPubSubResponseGetStatusCode(self.hresponse,
                                                p_value)
    if err != S_OK:
      self.logger.error("AzPubSubResponseGetStatusCode failed with: {0}".format(win32api.FormatMessage( err)))
      return None
    return p_value[0]

  def _get_response_sub_status_code(self):
    value = c_int(0)
    p_value = pointer(value)
    err = AzPubSub.AzPubSubResponseGetSubStatusCode(self.hresponse,
                                                p_value)
    if err != S_OK:
      self.logger.error("AzPubSubResponseGetSubStatusCode failed with: {0}".format(win32api.FormatMessage( err)))
      return None
    return p_value[0]

  def __del__(self):
    err = AzPubSub.AzPubSubResponseClose(self.hresponse)
    if err != S_OK:
      self.logger.error("AzPubSubResponseClose failed with: {0}".format(win32api.FormatMessage( err)))
    else:
      self.logger.info("AzPubSubResponse deleted successfully.")

  def __bool__(self):
    return self._status_code == SIMPLE_PRODUCER_SUCCESS_STATUS and \
    self.sub_status_code == SIMPLE_PRODUCER_SUCCESS_SUB_STATUS_OK


class GlobalResponse(AzPubSubResponse):
  """ Response details for the global client
  """  
  def __init__(self, status_code, message):
    # sometimes message_code can be None if the DLL doesn't return us an error code
    self._status_code = status_code
    self._message = message

  @property
  def message(self):
    return self._message

  @property
  def status_code(self):
    return self._status_code

  def __bool__(self):
    return self._status_code == 0


class AzPubSubClientBase():
  '''
  Base class for AzPubSub clients
  '''
  def __init__(self, endpoint, test_instance=False):
    self.logger = logging.getLogger('AzPubSubClient')
    self.endpoint = endpoint
    self.test_instance = test_instance


    self.hclient = None
    self.hconfig = None
    self.hproducer = None

    self.aps_config_type = None
    self.aps_security_type = None
    self.aps_connection_flags = None

    self.init_globals()

  def __del__(self):
    err = AzPubSub.AzPubSubFreeConfiguration(self.hconfig)

    if err != S_OK:
      self._log_error_hr("AzPubSubFreeConfiguration",  err)
    else:
      self.logger.info("AzPubSubClient config deleted successfully.")

    err = AzPubSub.AzPubSubClientClose(self.hclient)
    if err != S_OK:
      self._log_error_hr("AzPubSubClientClose",  err)
    else:
      self.logger.info("AzPubSubClient closed successfully.")

  def init_globals(self):
    if self.test_instance:
      self.logger.info("Initializing client with test globals")
      # This will force AzPubSub.dll to not load APPKIs
      self.aps_security_type = AZPUBSUB_SECURITY_TYPE.NONE
      # for pointing to localhost
      self.aps_connection_flags = AZPUBSUB_SECURITY_FLAGS.LOCAL
    else:
      self.logger.info("Initializing client with prod globals")
      self.aps_security_type = AZPUBSUB_SECURITY_TYPE.SSL
      self.aps_connection_flags = AZPUBSUB_SECURITY_FLAGS.NONE

  def init_config(self):
    self.hconfig = AzPubSub.AzPubSubCreateConfiguration(self.hclient, self.aps_config_type, AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES.NONE)

    if not self.hconfig:
      self.logger.error("Failed to initialize AzPubSubClient config: {0}".format(GetLastError()))
      return
    else:
      self.logger.info("AzPubSubClient config init successfully.")

  def init_client(self):
    self.hclient = AzPubSub.AzPubSubClientInitialize(p_logger_callback,
                                                     LPVOID())
    if not self.hclient:
      self.logger.error("Failed to initialize AzPubSubClient:{0}".format(GetLastError()))
    else:
      self.logger.info("AzPubSubClient init successfully.")

  def error_to_string(err_code):
    buffer_length = DWORD(0)
    p_buffer_length = LPDWORD(buffer_length)

    err = AzPubSub.AzPubSubErrorToString(err_code, NULL, p_buffer_length)
    if err != ERROR_MORE_DATA:
      self._log_error_hr("First call to AzPubSubErrorToString",  err)
      return None

    # increment buffer length to accommodate null char
    p_buffer_length.contents.value = p_buffer_length.contents.value + 1

    # allocate buffer
    p_msg_buffer = LPSTR(b'0'*buffer_length.value)
    pp_msg_buffer = pointer(p_msg_buffer)

    err = AzPubSub.AzPubSubErrorToString(err_code, pp_msg_buffer, p_buffer_length)
    if err != S_OK:
      self._log_error_hr("Second call to AzPubSubErrorToString",  err)
      return None

    return p_msg_buffer.value.decode("ascii")

  def add_config_key(self, key, value):
    err = S_OK
    if isinstance(value, str):
      err = AzPubSub.AzPubSubAddStringConfiguration(self.hconfig,
                                              LPCSTR(key.encode("ascii")),
                                              LPCSTR(value.encode("ascii")))
    elif isinstance(value, bool):
      err = AzPubSub.AzPubSubAddBooleanConfiguration(self.hconfig,
                                              LPCSTR(key.encode("ascii")),
                                              c_bool(value))
    elif isinstance(value, int):
      err = AzPubSub.AzPubSubAddIntegerConfiguration(self.hconfig,
                                               LPCSTR(key.encode("ascii")),
                                               c_int(value))
    else:
      self.logger.error("Unsupported type of value: {0}".format(type(value)))

    self._log_error_hr("AzPubSubAddConfiguration",  err)

  def get_config_value(self, key):
    buffer_length = DWORD(0)
    p_buffer_length = LPDWORD(buffer_length)

    err = AzPubSub.AzPubSubGetStringConfiguration(self.hconfig,
                                                 LPCSTR(key.encode("ascii")),
                                                 NULL,
                                                 buffer_length,
                                                 p_buffer_length)

    if err != ERROR_MORE_DATA:
      self._log_error_hr("First call to AzPubSubGetStringConfiguration",  err)
      return None

    # increment buffer length to accommodate null char
    p_buffer_length.contents.value = p_buffer_length.contents.value + 1

    # allocate buffer
    p_msg_buffer = LPSTR(b'0'*buffer_length.value)
    pp_msg_buffer = pointer(p_msg_buffer)

    err = AzPubSub.AzPubSubGetStringConfiguration(self.hconfig,
                                                 LPCSTR(key.encode("ascii")),
                                                 pp_msg_buffer,
                                                 buffer_length,
                                                 p_buffer_length)

    if err != S_OK:
      self._log_error_hr("Second call to AzPubSubGetStringConfiguration",  err)
      return None

    return p_msg_buffer.value.decode("ascii")

  def _log_error_hr(self,operation, err):
    if err != S_OK:
      self.logger.error("{0} failed with: {1}".format(operation,win32api.FormatMessage( err)))


class AzPubSubSimpleClient(AzPubSubClientBase):
  '''
  Class for synchronously publishing messages to AzPubSub via HTTP proxy. This has advantage of not maintaining
  persistent connection to AzPubSub endpoint, but send_message calls can take more time to complete.
  '''
  def __init__(self, endpoint, test_instance=False):
    super().__init__(endpoint, test_instance)
    # Required by simple producer
    self.aps_config_type = AZPUBSUB_CONFIGURATION_TYPE.SIMPLE
    self.init_client()
    self.init_config()

    self.init()

  def init(self):
    # set connection is not needed for simple producer
    self.open_simple_producer()

  def open_simple_producer(self):
    self.hproducer = AzPubSub.AzPubSubOpenSimpleProducer(self.hconfig,
                                        self.aps_security_type,
                                        NULL,
                                        NULL,
                                        LPCSTR(self.endpoint.encode("ascii")))
    if not self.hproducer:
      self.logger.error("AzPubSubOpenSimpleProducer failed with: {0}".format(win32api.FormatMessage(GetLastError())))
    else:
      self.logger.info("AzPubSubOpenSimpleProducer init successfully")

  def send_message(self, topic, msg, key=None):
    payload = msg.encode("ascii")
    hresponse = HRESPONSE(0)
    p_hresponse = pointer(hresponse)
    err = AzPubSub.AzPubSubSendMessageEx(self.hproducer,
                                        LPCSTR(topic.encode("ascii")),
                                        LPCSTR(key.encode("ascii")) if key != None else NULL,# Set key to NULL for default hash based partitioning
                                        NULL,
                                        LPCSTR(payload),
                                        len(payload),
                                        p_hresponse
                                        )
    self._log_error_hr("AzPubSubSendMessageEx",  err)

    if err != S_OK:
      return False

    response = SimpleResponse(hresponse)
    if not response:
      self.logger.error("send_message failed with {0},{1},{2}".format(response.status_code, response.sub_status_code, response.message))
    else:
      self.logger.info(
        "send_message success {0}".format(response.message))
    return response

class AzPubSubGlobalClient(AzPubSubClientBase):
  '''
  Class for synchronously publishing messages to AzPubSub. This maintais persistent connection to the
  AzPubSub endpoint, and can send headers as part of the message payload.
  '''
  def __init__(self, endpoint, topics=[], test_instance=False):
    super().__init__(endpoint, test_instance)
    self.opened_topics = {}
    self.aps_config_type = AZPUBSUB_CONFIGURATION_TYPE.GLOBAL
    self.init_client()
    self.init_config()

    self.init()

    for topic in topics:
      self.open_topic(topic)

  def __del__(self):
    for topic_name, hproducer_topic in self.opened_topics.items():
      err = AzPubSub.AzPubSubCloseProducerTopic(hproducer_topic)
      if err != S_OK:
        self._log_error_hr("AzPubSubCloseProducerTopic failure for topic {0}".format(topic_name), err)
      else:
        self.logger.info("AzPubSubCloseProducerTopic success for topic {0}".format(topic_name))

    if self.hproducer:
      err = AzPubSub.AzPubSubProducerClose(self.hproducer)
      if err != S_OK:
        self._log_error_hr("AzPubSubProducerClose", err)
      else:
        self.logger.info("AzPubSubProducerClose success.")

    super().__del__()

  def init(self):
    self.add_config_key("azpubsub.security.provider", "ApPki")
    self.connect()
    self.open_global_producer()

  def connect(self):
    err = AzPubSub.AzPubSubSetConnection(self.hconfig,
                                         NULL,
                                         NULL,
                                         LPCSTR(self.endpoint.encode("ascii")),
                                         NULL,
                                         self.aps_security_type,
                                         self.aps_connection_flags)
    self._log_error_hr("AzPubSubSetConnection",  err)

  def open_global_producer(self):
    self.hproducer = AzPubSub.AzPubSubOpenProducer(self.hconfig,
                                                   cast(NULL, AZPUBSUB_EVENT_CALLBACK),
                                                   cast(NULL, AZPUBSUB_METRICS_CALLBACK))
    if not self.hproducer:
      self.logger.error("AzPubSubOpenProducer failed with: {0}".format(win32api.FormatMessage(GetLastError())))
    else:
      self.logger.info("AzPubSubOpenProducer init successfully")

  def open_topic(self, topic):
    hproducer_topic = AzPubSub.AzPubSubOpenProducerTopic(self.hproducer,
                                                         LPCSTR(topic.encode("ascii")),
                                                         cast(NULL, HCONFIG),
                                                         cast(NULL, AZPUBSUB_TOPIC_PARTITIONER))
    if not hproducer_topic:
      self.logger.error("AzPubSubOpenProducerTopic for topic {0} failed with: {1}".format(topic, win32api.FormatMessage(GetLastError())))
      return False
    else:
      self.opened_topics[topic] = hproducer_topic
      self.logger.info("AzPubSubOpenProducerTopic success for topic {0}".format(topic))
      return True

  def send_message(self, topic, msg, key=None, headers=None, context=None):
    """
    param {context} must be encoded as an ascii string

    Encode your message with: your_context_str.encode('ascii')
    """
    payload = msg.encode("ascii")

    if topic not in self.opened_topics:
      res = self.open_topic(topic)
      if not res:
        self.logger.error("Failed to open topic {0} when sending message.".format(topic))
        return False

    err = AzPubSub.AzPubSubSendMessage(self.hproducer,
                                       self.opened_topics[topic],
                                       LPCSTR(key.encode("ascii")) if key else NULL,
                                       DWORD(len(key.encode("ascii"))) if key else 0,
                                       LPCSTR(payload),
                                       len(payload),
                                       -1,  # Set partition to -1 for default partitioning based on key
                                       0,
                                       encode_headers(headers) if headers else NULL,
                                       p_message_callback,
                                       LPCSTR(context) if context else NULL)
    if err != S_OK:
      self._log_error_hr("AzPubSubSendMessage",  err)
      return False

    return True

def encode_headers(headers):
  '''
  Encode the given headers dictionary into a memory block.

  Structure of the header block:
    - (1) DWORD representing the number of headers.
    - (2) Array of AZPUBSUB_MESSAGE_HEADER where the pointers in the structure points to
      strings and byte arrays after the array.
    - (3) Keys and Values packed at the end of the block.

  Inspired from MessageHeaderListMarshaler.cs and the GetHeaders method in Message.cpp.
  '''

  # Calculate buffer size
  buffer_size = sizeof(LPVOID) + len(headers) * sizeof(AZPUBSUB_MESSAGE_HEADER)
  for header_key, header_value in headers.items():
    buffer_size += len(header_key) + 1
    buffer_size += len(header_value) + 1

  # Allocate buffer
  buf = create_string_buffer(buffer_size)

  # Save number of headers
  buf[0] = len(headers)

  # Convert to LPVOID so pointer can be incremented
  buf = cast(buf, LPVOID)

  # Address of header array
  header_ptr = buf.value + sizeof(LPVOID)

  # Address of packed keys/values
  key_value_ptr = header_ptr + len(headers) * sizeof(AZPUBSUB_MESSAGE_HEADER)

  # Populate header array and key/value pairs
  for header_index, (header_key, header_value) in enumerate(headers.items()):
    # Header array pointer
    header = cast(header_ptr, POINTER(AZPUBSUB_MESSAGE_HEADER))

    # Save key pointer
    header.contents.pszKey = cast(LPVOID(key_value_ptr), LPCSTR)

    # Save key
    key = LPCSTR(header_key.encode("ascii"))
    keySize = len(header_key) + 1
    memmove(key_value_ptr, cast(key, PBYTE), keySize)
    key_value_ptr += keySize

    # Save value pointer
    header.contents.pValue = cast(LPVOID(key_value_ptr), PBYTE)

    # Save value size
    valSize = len(header_value) + 1
    header.contents.cbValue = DWORD(valSize)

    # Save value
    val = cast(header_value.encode("ascii"), PBYTE)
    memmove(key_value_ptr, val, valSize)
    key_value_ptr += valSize

    header_ptr += sizeof(AZPUBSUB_MESSAGE_HEADER)

  return buf

def decode_headers(encoded_headers):
  num_headers = cast(encoded_headers, POINTER(DWORD)).contents.value
  headers = {}
  header_structs = []
  key_value_ptr = encoded_headers.value + sizeof(LPVOID) + num_headers * sizeof(AZPUBSUB_MESSAGE_HEADER)
  for i in range(num_headers):
    header_struct_ptr = encoded_headers.value + sizeof(LPVOID) + i * sizeof(AZPUBSUB_MESSAGE_HEADER)
    header_struct = cast(header_struct_ptr, POINTER(AZPUBSUB_MESSAGE_HEADER))
    header_structs.append(header_struct.contents)

    key = cast(key_value_ptr, LPCSTR).value.decode('ascii')
    key_value_ptr += len(key) + 1

    value = cast(key_value_ptr, LPCSTR).value.decode('ascii')
    key_value_ptr += len(value) + 1

    headers[key] = value
  return num_headers, headers, header_structs

# SIG # Begin Windows Authenticode signature block
# MIInOAYJKoZIhvcNAQcCoIInKTCCJyUCAQExDzANBglghkgBZQMEAgEFADB5Bgor
# BgEEAYI3AgEEoGswaTA0BgorBgEEAYI3AgEeMCYCAwEAAAQQse8BENmB6EqSR2hd
# JGAGggIBAAIBAAIBAAIBAAIBADAxMA0GCWCGSAFlAwQCAQUABCDGzoTPoXkfuXDx
# JApu49pLAOw4JKQqYxYhnn0VMK4QCqCCEWUwggh3MIIHX6ADAgECAhM2AAABOXjG
# OfXldyfqAAEAAAE5MA0GCSqGSIb3DQEBCwUAMEExEzARBgoJkiaJk/IsZAEZFgNH
# QkwxEzARBgoJkiaJk/IsZAEZFgNBTUUxFTATBgNVBAMTDEFNRSBDUyBDQSAwMTAe
# Fw0yMDEwMjEyMDM5MDZaFw0yMTA5MTUyMTQzMDNaMCQxIjAgBgNVBAMTGU1pY3Jv
# c29mdCBBenVyZSBDb2RlIFNpZ24wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
# AoIBAQCvtf6RG9X1bFXLQOzLuA06k5gBhizLWQ3/m6nIKOwoNsu9N+s9yt+ZGRpb
# ZbDtBmtAeoi3c2XK9vf0x3sq32GWPPv+px6a7u55tQ9lq4evX6QNxPhrH++ltlUt
# siiVmV934/+F5B/71sJ1Nxr89OsExV1b5Ey7LiKkEwxpTRxlOyUXf4OiQvTDzG0I
# 7AseJ4RxOy23tLnh8268pkucY2PbSLFYoRIG1ZGNgchcprL+uiRLuCz4vZXfidQo
# Wus3ThY8+mYulD8AaQ5ZtnuwzSHtzxYm/g6OeSDsf4xFep0DYLA3zNiKO4CvmzNR
# jJbcg1Bm7OpDe/CSLSWG5aoqW+X5AgMBAAGjggWDMIIFfzApBgkrBgEEAYI3FQoE
# HDAaMAwGCisGAQQBgjdbAQEwCgYIKwYBBQUHAwMwPQYJKwYBBAGCNxUHBDAwLgYm
# KwYBBAGCNxUIhpDjDYTVtHiE8Ys+hZvdFs6dEoFgg93NZoaUjDICAWQCAQwwggJ2
# BggrBgEFBQcBAQSCAmgwggJkMGIGCCsGAQUFBzAChlZodHRwOi8vY3JsLm1pY3Jv
# c29mdC5jb20vcGtpaW5mcmEvQ2VydHMvQlkyUEtJQ1NDQTAxLkFNRS5HQkxfQU1F
# JTIwQ1MlMjBDQSUyMDAxKDEpLmNydDBSBggrBgEFBQcwAoZGaHR0cDovL2NybDEu
# YW1lLmdibC9haWEvQlkyUEtJQ1NDQTAxLkFNRS5HQkxfQU1FJTIwQ1MlMjBDQSUy
# MDAxKDEpLmNydDBSBggrBgEFBQcwAoZGaHR0cDovL2NybDIuYW1lLmdibC9haWEv
# QlkyUEtJQ1NDQTAxLkFNRS5HQkxfQU1FJTIwQ1MlMjBDQSUyMDAxKDEpLmNydDBS
# BggrBgEFBQcwAoZGaHR0cDovL2NybDMuYW1lLmdibC9haWEvQlkyUEtJQ1NDQTAx
# LkFNRS5HQkxfQU1FJTIwQ1MlMjBDQSUyMDAxKDEpLmNydDBSBggrBgEFBQcwAoZG
# aHR0cDovL2NybDQuYW1lLmdibC9haWEvQlkyUEtJQ1NDQTAxLkFNRS5HQkxfQU1F
# JTIwQ1MlMjBDQSUyMDAxKDEpLmNydDCBrQYIKwYBBQUHMAKGgaBsZGFwOi8vL0NO
# PUFNRSUyMENTJTIwQ0ElMjAwMSxDTj1BSUEsQ049UHVibGljJTIwS2V5JTIwU2Vy
# dmljZXMsQ049U2VydmljZXMsQ049Q29uZmlndXJhdGlvbixEQz1BTUUsREM9R0JM
# P2NBQ2VydGlmaWNhdGU/YmFzZT9vYmplY3RDbGFzcz1jZXJ0aWZpY2F0aW9uQXV0
# aG9yaXR5MB0GA1UdDgQWBBRQasfWFuGWZ4TjHj7E0G+JYLldgzAOBgNVHQ8BAf8E
# BAMCB4AwUAYDVR0RBEkwR6RFMEMxKTAnBgNVBAsTIE1pY3Jvc29mdCBPcGVyYXRp
# b25zIFB1ZXJ0byBSaWNvMRYwFAYDVQQFEw0yMzYxNjcrNDYyNTE2MIIB1AYDVR0f
# BIIByzCCAccwggHDoIIBv6CCAbuGPGh0dHA6Ly9jcmwubWljcm9zb2Z0LmNvbS9w
# a2lpbmZyYS9DUkwvQU1FJTIwQ1MlMjBDQSUyMDAxLmNybIYuaHR0cDovL2NybDEu
# YW1lLmdibC9jcmwvQU1FJTIwQ1MlMjBDQSUyMDAxLmNybIYuaHR0cDovL2NybDIu
# YW1lLmdibC9jcmwvQU1FJTIwQ1MlMjBDQSUyMDAxLmNybIYuaHR0cDovL2NybDMu
# YW1lLmdibC9jcmwvQU1FJTIwQ1MlMjBDQSUyMDAxLmNybIYuaHR0cDovL2NybDQu
# YW1lLmdibC9jcmwvQU1FJTIwQ1MlMjBDQSUyMDAxLmNybIaBumxkYXA6Ly8vQ049
# QU1FJTIwQ1MlMjBDQSUyMDAxLENOPUJZMlBLSUNTQ0EwMSxDTj1DRFAsQ049UHVi
# bGljJTIwS2V5JTIwU2VydmljZXMsQ049U2VydmljZXMsQ049Q29uZmlndXJhdGlv
# bixEQz1BTUUsREM9R0JMP2NlcnRpZmljYXRlUmV2b2NhdGlvbkxpc3Q/YmFzZT9v
# YmplY3RDbGFzcz1jUkxEaXN0cmlidXRpb25Qb2ludDAfBgNVHSMEGDAWgBQbZqIZ
# /JvrpdqEjxiY6RCkw3uSvTAfBgNVHSUEGDAWBgorBgEEAYI3WwEBBggrBgEFBQcD
# AzANBgkqhkiG9w0BAQsFAAOCAQEArFNMfAJStrd/3V4hInTdjEo/CLYAY8YX/foG
# Amyk6NrjEx3uFN0sJmR3qR0iBggS3SBiUi4oZ+Xk8+DjVnnJFn9Fhmu/kB2wT4ZK
# jjjZeWROPcTsUnRgs1+OhKTWbX2Eng8oH3Cq0qR9LaOT/ES5Ejd98S1jq6WZ8B8K
# dNHg0d+VGAtwts+E3uu8MkUM5rUukmPHW7BC8ttmgKeXZiIiLV4T1KzxBMMNg0lY
# 7iFbQ5fkj5hLa1E0WvsGMcMGOMwRUVwVwl6F8OL8aUY5i7tpAuz54XVS4W1grPyT
# JDae1qB19H5JvqTwPPNm30JrFGpR/X/SGQhROsoD4V1tvCJ8tDCCCOYwggbOoAMC
# AQICEx8AAAAUtMUfxvKAvnEAAAAAABQwDQYJKoZIhvcNAQELBQAwPDETMBEGCgmS
# JomT8ixkARkWA0dCTDETMBEGCgmSJomT8ixkARkWA0FNRTEQMA4GA1UEAxMHYW1l
# cm9vdDAeFw0xNjA5MTUyMTMzMDNaFw0yMTA5MTUyMTQzMDNaMEExEzARBgoJkiaJ
# k/IsZAEZFgNHQkwxEzARBgoJkiaJk/IsZAEZFgNBTUUxFTATBgNVBAMTDEFNRSBD
# UyBDQSAwMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANVXgQLW+frQ
# 9xuAud03zSTcZmH84YlyrSkM0hsbmr+utG00tVRHgw40pxYbJp5W+hpDwnmJgicF
# oGRrPt6FifMmnd//1aD/fW1xvGs80yZk9jxTNcisVF1CYIuyPctwuJZfwE3wcGxh
# kVw/tj3ZHZVacSls3jRD1cGwrcVo1IR6+hHMvUejtt4/tv0UmUoH82HLQ8w1oTX9
# D7xj35Zt9T0pOPqM3Gt9+/zs7tPp2gyoOYv8xR4X0iWZKuXTzxugvMA63YsB4ehu
# SBqzHdkF55rxH47aT6hPhvDHlm7M2lsZcRI0CUAujwcJ/vELeFapXNGpt2d3wcPJ
# M0bpzrPDJ/8CAwEAAaOCBNowggTWMBAGCSsGAQQBgjcVAQQDAgEBMCMGCSsGAQQB
# gjcVAgQWBBSR/DPOQp72k+bifVTXCBi7uNdxZTAdBgNVHQ4EFgQUG2aiGfyb66Xa
# hI8YmOkQpMN7kr0wggEEBgNVHSUEgfwwgfkGBysGAQUCAwUGCCsGAQUFBwMBBggr
# BgEFBQcDAgYKKwYBBAGCNxQCAQYJKwYBBAGCNxUGBgorBgEEAYI3CgMMBgkrBgEE
# AYI3FQYGCCsGAQUFBwMJBggrBgEFBQgCAgYKKwYBBAGCN0ABAQYLKwYBBAGCNwoD
# BAEGCisGAQQBgjcKAwQGCSsGAQQBgjcVBQYKKwYBBAGCNxQCAgYKKwYBBAGCNxQC
# AwYIKwYBBQUHAwMGCisGAQQBgjdbAQEGCisGAQQBgjdbAgEGCisGAQQBgjdbAwEG
# CisGAQQBgjdbBQEGCisGAQQBgjdbBAEGCisGAQQBgjdbBAIwGQYJKwYBBAGCNxQC
# BAweCgBTAHUAYgBDAEEwCwYDVR0PBAQDAgGGMBIGA1UdEwEB/wQIMAYBAf8CAQAw
# HwYDVR0jBBgwFoAUKV5RXmSuNLnrrJwNp4x1AdEJCygwggFoBgNVHR8EggFfMIIB
# WzCCAVegggFToIIBT4YjaHR0cDovL2NybDEuYW1lLmdibC9jcmwvYW1lcm9vdC5j
# cmyGMWh0dHA6Ly9jcmwubWljcm9zb2Z0LmNvbS9wa2lpbmZyYS9jcmwvYW1lcm9v
# dC5jcmyGI2h0dHA6Ly9jcmwyLmFtZS5nYmwvY3JsL2FtZXJvb3QuY3JshiNodHRw
# Oi8vY3JsMy5hbWUuZ2JsL2NybC9hbWVyb290LmNybIaBqmxkYXA6Ly8vQ049YW1l
# cm9vdCxDTj1BTUVST09ULENOPUNEUCxDTj1QdWJsaWMlMjBLZXklMjBTZXJ2aWNl
# cyxDTj1TZXJ2aWNlcyxDTj1Db25maWd1cmF0aW9uLERDPUFNRSxEQz1HQkw/Y2Vy
# dGlmaWNhdGVSZXZvY2F0aW9uTGlzdD9iYXNlP29iamVjdENsYXNzPWNSTERpc3Ry
# aWJ1dGlvblBvaW50MIIBqwYIKwYBBQUHAQEEggGdMIIBmTA3BggrBgEFBQcwAoYr
# aHR0cDovL2NybDEuYW1lLmdibC9haWEvQU1FUk9PVF9hbWVyb290LmNydDBHBggr
# BgEFBQcwAoY7aHR0cDovL2NybC5taWNyb3NvZnQuY29tL3BraWluZnJhL2NlcnRz
# L0FNRVJPT1RfYW1lcm9vdC5jcnQwNwYIKwYBBQUHMAKGK2h0dHA6Ly9jcmwyLmFt
# ZS5nYmwvYWlhL0FNRVJPT1RfYW1lcm9vdC5jcnQwNwYIKwYBBQUHMAKGK2h0dHA6
# Ly9jcmwzLmFtZS5nYmwvYWlhL0FNRVJPT1RfYW1lcm9vdC5jcnQwgaIGCCsGAQUF
# BzAChoGVbGRhcDovLy9DTj1hbWVyb290LENOPUFJQSxDTj1QdWJsaWMlMjBLZXkl
# MjBTZXJ2aWNlcyxDTj1TZXJ2aWNlcyxDTj1Db25maWd1cmF0aW9uLERDPUFNRSxE
# Qz1HQkw/Y0FDZXJ0aWZpY2F0ZT9iYXNlP29iamVjdENsYXNzPWNlcnRpZmljYXRp
# b25BdXRob3JpdHkwDQYJKoZIhvcNAQELBQADggIBACi3Soaajx+kAWjNwgDqkIvK
# AOFkHmS1t0DlzZlpu1ANNfA0BGtck6hEG7g+TpUdVrvxdvPQ5lzU3bGTOBkyhGmX
# oSIlWjKC7xCbbuYegk8n1qj3rTcjiakdbBqqHdF8J+fxv83E2EsZ+StzfCnZXA62
# QCMn6t8mhCWBxpwPXif39Ua32yYHqP0QISAnLTjjcH6bAV3IIk7k5pQ/5NA6qIL8
# yYD6vRjpCMl/3cZOyJD81/5+POLNMx0eCClOfFNxtaD0kJmeThwL4B2hAEpHTeRN
# tB8ib+cze3bvkGNPHyPlSHIuqWoC31x2Gk192SfzFDPV1PqFOcuKjC8049SSBtC1
# X7hyvMqAe4dop8k3u25+odhvDcWdNmimdMWvp/yZ6FyjbGlTxtUqE7iLTLF1eaUL
# SEobAap16hY2N2yTJTISKHzHI4rjsEQlvqa2fj6GLxNj/jC+4LNy+uRmfQXShd30
# lt075qTroz0Nt680pXvVhsRSdNnzW2hfQu2xuOLg8zKGVOD/rr0GgeyhODjKgL2G
# Hxctbb9XaVSDf6ocdB//aDYjiabmWd/WYmy7fQ127KuasMh5nSV2orMcAed8CbIV
# I3NYu+sahT1DRm/BGUN2hSpdsPQeO73wYvp1N7DdLaZyz7XsOCx1quCwQ+bojWVQ
# TmKLGegSoUpZNfmP9MtSMYIVKTCCFSUCAQEwWDBBMRMwEQYKCZImiZPyLGQBGRYD
# R0JMMRMwEQYKCZImiZPyLGQBGRYDQU1FMRUwEwYDVQQDEwxBTUUgQ1MgQ0EgMDEC
# EzYAAAE5eMY59eV3J+oAAQAAATkwDQYJYIZIAWUDBAIBBQCgga4wGQYJKoZIhvcN
# AQkDMQwGCisGAQQBgjcCAQQwHAYKKwYBBAGCNwIBCzEOMAwGCisGAQQBgjcCARUw
# LwYJKoZIhvcNAQkEMSIEIHh1ec7r7G82UIoeq1ZREa/mBkweS5ReIRdDkbdX2Qdo
# MEIGCisGAQQBgjcCAQwxNDAyoBSAEgBNAGkAYwByAG8AcwBvAGYAdKEagBhodHRw
# Oi8vd3d3Lm1pY3Jvc29mdC5jb20wDQYJKoZIhvcNAQEBBQAEggEASs/5OzTY6QeF
# 6M5ZDW1ythJWMkZTU5VPjN0HLQnKX/zSXCoHxygk1vYZ/6oPz76pdSlaa4E1+ghC
# pm7eU2qjRx8D8lINdVlmJXjCIxOLpF9CGcrpg7eAaR69vSdRKw4FGBOhUlSQMgZY
# nbXlW54RNuOZQtKG0gNs2yYljjG/KKEQypFwGfrbiSzZMqRyb5ToJ/WyPX205xbZ
# SmN+sP3gHfz7ueeEJFnUhuuoOs6ACVD8ZWqh34ihTtz2Lp4CVCB0rf3mZ0Iec3W5
# dXUk6cp0twnNKqWNuFafrw5kYdCST5Y9JW5NBB9LuyOStvlzTsSBGue3BJuD99rT
# 6T4N1PLvgKGCEvEwghLtBgorBgEEAYI3AwMBMYIS3TCCEtkGCSqGSIb3DQEHAqCC
# EsowghLGAgEDMQ8wDQYJYIZIAWUDBAIBBQAwggFVBgsqhkiG9w0BCRABBKCCAUQE
# ggFAMIIBPAIBAQYKKwYBBAGEWQoDATAxMA0GCWCGSAFlAwQCAQUABCC9BzW+XOFp
# gjbu1dMhBBjB6P/Kb5IQlIzxrVXS8hKB8AIGYK67cs90GBMyMDIxMDYwODE4NDYw
# OC41NTVaMASAAgH0oIHUpIHRMIHOMQswCQYDVQQGEwJVUzETMBEGA1UECBMKV2Fz
# aGluZ3RvbjEQMA4GA1UEBxMHUmVkbW9uZDEeMBwGA1UEChMVTWljcm9zb2Z0IENv
# cnBvcmF0aW9uMSkwJwYDVQQLEyBNaWNyb3NvZnQgT3BlcmF0aW9ucyBQdWVydG8g
# UmljbzEmMCQGA1UECxMdVGhhbGVzIFRTUyBFU046Rjg3QS1FMzc0LUQ3QjkxJTAj
# BgNVBAMTHE1pY3Jvc29mdCBUaW1lLVN0YW1wIFNlcnZpY2Wggg5EMIIE9TCCA92g
# AwIBAgITMwAAAWOLZMbJhZZldgAAAAABYzANBgkqhkiG9w0BAQsFADB8MQswCQYD
# VQQGEwJVUzETMBEGA1UECBMKV2FzaGluZ3RvbjEQMA4GA1UEBxMHUmVkbW9uZDEe
# MBwGA1UEChMVTWljcm9zb2Z0IENvcnBvcmF0aW9uMSYwJAYDVQQDEx1NaWNyb3Nv
# ZnQgVGltZS1TdGFtcCBQQ0EgMjAxMDAeFw0yMTAxMTQxOTAyMjNaFw0yMjA0MTEx
# OTAyMjNaMIHOMQswCQYDVQQGEwJVUzETMBEGA1UECBMKV2FzaGluZ3RvbjEQMA4G
# A1UEBxMHUmVkbW9uZDEeMBwGA1UEChMVTWljcm9zb2Z0IENvcnBvcmF0aW9uMSkw
# JwYDVQQLEyBNaWNyb3NvZnQgT3BlcmF0aW9ucyBQdWVydG8gUmljbzEmMCQGA1UE
# CxMdVGhhbGVzIFRTUyBFU046Rjg3QS1FMzc0LUQ3QjkxJTAjBgNVBAMTHE1pY3Jv
# c29mdCBUaW1lLVN0YW1wIFNlcnZpY2UwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
# ggEKAoIBAQCtcRf2Ep3JdGKS/6jdhZ38I39IvGvguBC8+VGctyTlOeFx/qx79cty
# 4CmBrt8K7TUhGOB8+c+j0ZRCb7+2itrZu1rxUnrO4ixYUNPA1eIVcecxepZebjdr
# YtnyTWeiQ4zElWLmP8GmHTRaOzeJfMfO/9UkyKG9zw4mqgKBGdYRG5rka+OBCj/9
# 0Q4KPwGNKNNcwBeJOR78q389NxmiSGehCCIG2GxOhNi19nCWfet2jWD2S2FWzZ07
# 4ju6dnhh7WgJJ9PEK81vac9Whgk1JQy0VC5zIkFSzYoGlNb/Dk87+2pQCJ05UXxS
# 7zyFdCSdkj6vsFS8TxoYlbMBK1/fP7M1AgMBAAGjggEbMIIBFzAdBgNVHQ4EFgQU
# CTXK8XZyZ+4/MVqfRseQPtffPSkwHwYDVR0jBBgwFoAU1WM6XIoxkPNDe3xGG8Uz
# aFqFbVUwVgYDVR0fBE8wTTBLoEmgR4ZFaHR0cDovL2NybC5taWNyb3NvZnQuY29t
# L3BraS9jcmwvcHJvZHVjdHMvTWljVGltU3RhUENBXzIwMTAtMDctMDEuY3JsMFoG
# CCsGAQUFBwEBBE4wTDBKBggrBgEFBQcwAoY+aHR0cDovL3d3dy5taWNyb3NvZnQu
# Y29tL3BraS9jZXJ0cy9NaWNUaW1TdGFQQ0FfMjAxMC0wNy0wMS5jcnQwDAYDVR0T
# AQH/BAIwADATBgNVHSUEDDAKBggrBgEFBQcDCDANBgkqhkiG9w0BAQsFAAOCAQEA
# AohBggfXjuJTzo4yAmH7E6mpvSKnUbTI9tFAQVS4bn7z/cb5aCPC2fcDj6uLAqCU
# nYTC2sFFmXeu7xZTP4gT/u15KtdPU2nkEhODXPbnjNeX5RL2qOGbcxqFk3MaQvmp
# WGNJFRiI+ksQUsZwpKGXrE+OFlSEwUC/+Nz5h8VQBQ9AtXA882uZ79Qkog752eKj
# caT+mn/SGHymyQeGycQaudhWVUKkeHQOjWux+LE4YdQGP6mHOpM5kqYVLxMwqucT
# 2fPk5bKDTWWM+kwEeqp3n09g/9w7J+15jvsDYyIugBFkCR2qsAe0eTlju0Ce6dO0
# Zf+E75DTM72ZfAQUn1+2IzCCBnEwggRZoAMCAQICCmEJgSoAAAAAAAIwDQYJKoZI
# hvcNAQELBQAwgYgxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpXYXNoaW5ndG9uMRAw
# DgYDVQQHEwdSZWRtb25kMR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24x
# MjAwBgNVBAMTKU1pY3Jvc29mdCBSb290IENlcnRpZmljYXRlIEF1dGhvcml0eSAy
# MDEwMB4XDTEwMDcwMTIxMzY1NVoXDTI1MDcwMTIxNDY1NVowfDELMAkGA1UEBhMC
# VVMxEzARBgNVBAgTCldhc2hpbmd0b24xEDAOBgNVBAcTB1JlZG1vbmQxHjAcBgNV
# BAoTFU1pY3Jvc29mdCBDb3Jwb3JhdGlvbjEmMCQGA1UEAxMdTWljcm9zb2Z0IFRp
# bWUtU3RhbXAgUENBIDIwMTAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
# AQCpHQ28dxGKOiDs/BOX9fp/aZRrdFQQ1aUKAIKF++18aEssX8XD5WHCdrc+Zitb
# 8BVTJwQxH0EbGpUdzgkTjnxhMFmxMEQP8WCIhFRDDNdNuDgIs0Ldk6zWczBXJoKj
# RQ3Q6vVHgc2/JGAyWGBG8lhHhjKEHnRhZ5FfgVSxz5NMksHEpl3RYRNuKMYa+YaA
# u99h/EbBJx0kZxJyGiGKr0tkiVBisV39dx898Fd1rL2KQk1AUdEPnAY+Z3/1ZsAD
# lkR+79BL/W7lmsqxqPJ6Kgox8NpOBpG2iAg16HgcsOmZzTznL0S6p/TcZL2kAcEg
# CZN4zfy8wMlEXV4WnAEFTyJNAgMBAAGjggHmMIIB4jAQBgkrBgEEAYI3FQEEAwIB
# ADAdBgNVHQ4EFgQU1WM6XIoxkPNDe3xGG8UzaFqFbVUwGQYJKwYBBAGCNxQCBAwe
# CgBTAHUAYgBDAEEwCwYDVR0PBAQDAgGGMA8GA1UdEwEB/wQFMAMBAf8wHwYDVR0j
# BBgwFoAU1fZWy4/oolxiaNE9lJBb186aGMQwVgYDVR0fBE8wTTBLoEmgR4ZFaHR0
# cDovL2NybC5taWNyb3NvZnQuY29tL3BraS9jcmwvcHJvZHVjdHMvTWljUm9vQ2Vy
# QXV0XzIwMTAtMDYtMjMuY3JsMFoGCCsGAQUFBwEBBE4wTDBKBggrBgEFBQcwAoY+
# aHR0cDovL3d3dy5taWNyb3NvZnQuY29tL3BraS9jZXJ0cy9NaWNSb29DZXJBdXRf
# MjAxMC0wNi0yMy5jcnQwgaAGA1UdIAEB/wSBlTCBkjCBjwYJKwYBBAGCNy4DMIGB
# MD0GCCsGAQUFBwIBFjFodHRwOi8vd3d3Lm1pY3Jvc29mdC5jb20vUEtJL2RvY3Mv
# Q1BTL2RlZmF1bHQuaHRtMEAGCCsGAQUFBwICMDQeMiAdAEwAZQBnAGEAbABfAFAA
# bwBsAGkAYwB5AF8AUwB0AGEAdABlAG0AZQBuAHQALiAdMA0GCSqGSIb3DQEBCwUA
# A4ICAQAH5ohRDeLG4Jg/gXEDPZ2joSFvs+umzPUxvs8F4qn++ldtGTCzwsVmyWrf
# 9efweL3HqJ4l4/m87WtUVwgrUYJEEvu5U4zM9GASinbMQEBBm9xcF/9c+V4XNZgk
# Vkt070IQyK+/f8Z/8jd9Wj8c8pl5SpFSAK84Dxf1L3mBZdmptWvkx872ynoAb0sw
# RCQiPM/tA6WWj1kpvLb9BOFwnzJKJ/1Vry/+tuWOM7tiX5rbV0Dp8c6ZZpCM/2pi
# f93FSguRJuI57BlKcWOdeyFtw5yjojz6f32WapB4pm3S4Zz5Hfw42JT0xqUKloak
# vZ4argRCg7i1gJsiOCC1JeVk7Pf0v35jWSUPei45V3aicaoGig+JFrphpxHLmtgO
# R5qAxdDNp9DvfYPw4TtxCd9ddJgiCGHasFAeb73x4QDf5zEHpJM692VHeOj4qEir
# 995yfmFrb3epgcunCaw5u+zGy9iCtHLNHfS4hQEegPsbiSpUObJb2sgNVZl6h3M7
# COaYLeqN4DMuEin1wC9UJyH3yKxO2ii4sanblrKnQqLJzxlBTeCG+SqaoxFmMNO7
# dDJL32N79ZmKLxvHIa9Zta7cRDyXUHHXodLFVeNp3lfB0d4wwP3M5k37Db9dT+md
# Hhk4L7zPWAUu7w2gUDXa7wknHNWzfjUeCLraNtvTX4/edIhJEqGCAtIwggI7AgEB
# MIH8oYHUpIHRMIHOMQswCQYDVQQGEwJVUzETMBEGA1UECBMKV2FzaGluZ3RvbjEQ
# MA4GA1UEBxMHUmVkbW9uZDEeMBwGA1UEChMVTWljcm9zb2Z0IENvcnBvcmF0aW9u
# MSkwJwYDVQQLEyBNaWNyb3NvZnQgT3BlcmF0aW9ucyBQdWVydG8gUmljbzEmMCQG
# A1UECxMdVGhhbGVzIFRTUyBFU046Rjg3QS1FMzc0LUQ3QjkxJTAjBgNVBAMTHE1p
# Y3Jvc29mdCBUaW1lLVN0YW1wIFNlcnZpY2WiIwoBATAHBgUrDgMCGgMVAO0sYB7d
# Sd0qk00qsy3KzBmUAWHvoIGDMIGApH4wfDELMAkGA1UEBhMCVVMxEzARBgNVBAgT
# Cldhc2hpbmd0b24xEDAOBgNVBAcTB1JlZG1vbmQxHjAcBgNVBAoTFU1pY3Jvc29m
# dCBDb3Jwb3JhdGlvbjEmMCQGA1UEAxMdTWljcm9zb2Z0IFRpbWUtU3RhbXAgUENB
# IDIwMTAwDQYJKoZIhvcNAQEFBQACBQDkabRiMCIYDzIwMjEwNjA4MTMxODI2WhgP
# MjAyMTA2MDkxMzE4MjZaMHcwPQYKKwYBBAGEWQoEATEvMC0wCgIFAORptGICAQAw
# CgIBAAICEkwCAf8wBwIBAAICESMwCgIFAORrBeICAQAwNgYKKwYBBAGEWQoEAjEo
# MCYwDAYKKwYBBAGEWQoDAqAKMAgCAQACAwehIKEKMAgCAQACAwGGoDANBgkqhkiG
# 9w0BAQUFAAOBgQAGOwVrvcV63XeBjCS63ZmKwDeUlCyzJW1G8v3u3hNU+JbOpRk1
# i1pePFQlEfH4bzLU2K+V8cKyj1c4wZGP8LuJrfltaAK4vckT7xabAI1vdletr604
# smDBK8nyQZaWGKgEOhKmRNG3bbbs91ad0ZCTqka+4gP0JkO2drmwVggkWzGCAw0w
# ggMJAgEBMIGTMHwxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpXYXNoaW5ndG9uMRAw
# DgYDVQQHEwdSZWRtb25kMR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24x
# JjAkBgNVBAMTHU1pY3Jvc29mdCBUaW1lLVN0YW1wIFBDQSAyMDEwAhMzAAABY4tk
# xsmFlmV2AAAAAAFjMA0GCWCGSAFlAwQCAQUAoIIBSjAaBgkqhkiG9w0BCQMxDQYL
# KoZIhvcNAQkQAQQwLwYJKoZIhvcNAQkEMSIEII5VC9t+Ur7jxfG0JLpA6Mlrj2Jw
# 2bqwgyBHN6UoTHIeMIH6BgsqhkiG9w0BCRACLzGB6jCB5zCB5DCBvQQgnFndlx2h
# Y6EopCm4uMvQGASKSwcvUW9ep7NRqxH0I2owgZgwgYCkfjB8MQswCQYDVQQGEwJV
# UzETMBEGA1UECBMKV2FzaGluZ3RvbjEQMA4GA1UEBxMHUmVkbW9uZDEeMBwGA1UE
# ChMVTWljcm9zb2Z0IENvcnBvcmF0aW9uMSYwJAYDVQQDEx1NaWNyb3NvZnQgVGlt
# ZS1TdGFtcCBQQ0EgMjAxMAITMwAAAWOLZMbJhZZldgAAAAABYzAiBCBU6Zv0zQDk
# EOSy+evQV71CpuzXWnZQssAUVYnUj92LHjANBgkqhkiG9w0BAQsFAASCAQAHjtsL
# mWWZBS2PIq7ObABSni/juVTuyuCwHMTT3SfhxQAp68B/UTEq3QYJ38TJ5QqslH/n
# Fz2FmN1eHa87yzBkb7lYpY6jTIKnLZu3ASlaL3+djpnPRla3Eh06Pr4FkKUSDzG+
# lAOUaK+0xC6RaaiE9y6L9w4OQj4Y7fK9DBRig8Ew/qpSjfZzscEWwXxCwKbJBQjW
# egj0mLmU/8F5dCZpmaOt5ibpeLsoDiRDzLvKDC5eOjqguVEBcFfH7uxBys0it2Li
# 3hA11cKptQgGHv3VTw+eBwn+oQf+FcNDIRLdCZgBkCzjRacuY6zrdX9imLDyfl3X
# yHmuDmxranYy0VZe
# SIG # End Windows Authenticode signature block