/**
@file

@brief AzPubSub Native Layer public header file. 
*/

#ifndef _AZ_PUBSUB_
#define _AZ_PUBSUB_

#if defined (_MSC_VER) && (_MSC_VER >= 1020)
#pragma once
#endif

#if defined (_MSC_VER)
#if ( _MSC_VER >= 1900 )
#ifdef __cplusplus

#ifdef AZPUBSUB_EXPORTS
#define AZPUBSUB_API __declspec(dllexport)
#else
#define AZPUBSUB_API __declspec(dllimport)
#endif

#ifdef __cplusplus
extern "C" {
#endif

    #define AZPUBSUB_CLIENT_VERSION_1    1
    #define AZPUBSUB_CLIENT_VERSION AZPUBSUB_CLIENT_VERSION_1

    typedef struct _HAZPUBSUBCLIENT const *HCLIENT;
    typedef struct _HAZPUBSUBCONFIG const *HCONFIG;
    typedef struct _HAZPUBSUBCONSUMER const *HCONSUMER;
    typedef struct _HAZPUBSUBPRODUCER const *HPRODUCER;
    typedef struct _HAZPUBSUBMESSAGE const *HMESSAGE;
    typedef struct _HAZPUBSUBRESPONSE const* HRESPONSE;
    typedef struct _HAZPUBSUBMETADATA const* HMETADATA;
    typedef struct _HAZPUBSUBCONNECTION const* HCONNECTION;
    typedef struct _HAZPUBSUBHPRODUCERTOPIC const* HPRODUCERTOPIC;
    typedef struct _HAZPUBSUBTOPICPARTITION const* HTOPICPARTITION;
    typedef struct _HAZPUBSUBTOPICPARTITIONLIST const* HTOPICPARTITIONLIST;
    typedef struct _HAZPUBSUBBATCH const* HBATCH;
    typedef struct _HTOKENOPERATION const* HTOKENOPERATION;

    /**
    Enum defining the different log levels.
    */
    typedef enum _AZPUBSUB_LOG_LEVEL {
        AZPUBSUB_LOG_LEVEL_EMERGENCY,
        AZPUBSUB_LOG_LEVEL_ALERT,
        AZPUBSUB_LOG_LEVEL_CRITICAL,
        AZPUBSUB_LOG_LEVEL_ERROR,
        AZPUBSUB_LOG_LEVEL_WARNING,
        AZPUBSUB_LOG_LEVEL_NOTICE,
        AZPUBSUB_LOG_LEVEL_INFO,
        AZPUBSUB_LOG_LEVEL_DEBUG,
    } AZPUBSUB_LOG_LEVEL;
    
    /**
    Enum defining the different event types.
    */
    typedef enum _AZPUBSUB_EVENT_TYPE {
        AZPUBSUB_EVENT_ERROR,
        AZPUBSUB_EVENT_STATS,
        AZPUBSUB_EVENT_LOG,
        AZPUBSUB_EVENT_THROTTLE,
    } AZPUBSUB_EVENT_TYPE;

    /**
    Enum defining the different security types.
    */
    typedef enum _AZPUBSUB_SECURITY_TYPE {
        AZPUBSUB_SECURITY_TYPE_NONE,
        AZPUBSUB_SECURITY_TYPE_SSL,
        AZPUBSUB_SECURITY_TYPE_TOKEN,
    } AZPUBSUB_SECURITY_TYPE, *PAZPUBSUB_SECURITY_TYPE;

    /**
    Enum defining the different connection flags.
    */
    typedef enum _AZPUBSUB_CONNECTION_FLAGS {
        AZPUBSUB_CONNECTION_FLAG_NONE = 0,
        AZPUBSUB_CONNECTION_FLAG_LOCAL = 1,
        AZPUBSUB_CONNECTION_FLAG_REGIONAL = 2
    } AZPUBSUB_CONNECTION_FLAGS, *PAZPUBSUB_CONNECTION_FLAGS;

    /**
    Enum defining the different rebalance modes.
    */
    typedef enum _AZPUBSUB_REBALANCE_MODE {
        AZPUBSUB_REBALANCE_MODE_NONE,
        AZPUBSUB_REBALANCE_MODE_ASSIGN,
        AZPUBSUB_REBALANCE_MODE_REVOKE,
        AZPUBSUB_REBALANCE_MODE_ERROR
    } AZPUBSUB_REBALANCE_MODE, *PAZPUBSUB_REBALANCE_MODE;

    /**
    The log callback function for a AzPubSub client. 

    @param level The log level of the message 
    @param lpMessage The message 
    @param context The context of the message
    */
    typedef void
    (WINAPI *AZPUBSUB_LOG_CALLBACK)(
        AZPUBSUB_LOG_LEVEL level,
        LPCSTR lpMessage,
        PVOID context);

    /**
    A message callback function.

    @param hMessage The message
    @param context The context of the message
    */
    typedef void
    (WINAPI *AZPUBSUB_MESSAGE_CALLBACK)(
        HMESSAGE hMessage,
        PVOID context);

    /**
    This is the rebalance callback.

    @param The consumer inciting the rebalance callback 
    @param mode The mode of the rebalance callback 
    @param hList The partitions to rebalance 
    */
    typedef void
    (WINAPI *AZPUBSUB_REBALANCE_CALLBACK)(
        HCONSUMER hConsumer,
        AZPUBSUB_REBALANCE_MODE mode,
        HTOPICPARTITIONLIST hList);

    /**
    The commit callback.

    @param error The error from the commit 
    @param the partition list
    */
    typedef void
    (WINAPI *AZPUBSUB_COMMIT_CALLBACK)(
        DWORD error,
        HTOPICPARTITIONLIST hList);

    typedef void
    (WINAPI *AZPUBSUB_EVENT_CALLBACK)(
        AZPUBSUB_EVENT_TYPE type,
        DWORD ErrorCode,
        LPCSTR lpszErrorMessage,
        AZPUBSUB_LOG_LEVEL level,
        LPCSTR lpszFacility,
        LPCSTR lpszMessage,
        DWORD throttle,
        LPCSTR lpszName,
        DWORD id);

    typedef void
    (WINAPI *AZPUBSUB_METRICS_CALLBACK)(
        LPCSTR lpszName,
        ULONGLONG value);

    typedef DWORD
    (WINAPI *AZPUBSUB_SET_TOKEN_CALLBACK)(
        HTOKENOPERATION hOperation,
        LPCSTR lpszToken);

    typedef DWORD
    (WINAPI *AZPUBSUB_TOKEN_REFRESH_CALLBACK)(
        AZPUBSUB_SET_TOKEN_CALLBACK SetToken,
        ULONGLONG* pRefresh,
        HTOKENOPERATION hOperation);

    typedef DWORD
    (WINAPI *AZPUBSUB_TOPIC_PARTITIONER)(
        HPRODUCERTOPIC hProducerTopic,
        const PBYTE key,
        DWORD cbKey,
        DWORD partitions,
        void* context);

    /**
    An enum defining the configuration types.
    */
    typedef enum _AZPUBSUB_CONFIGURATION_TYPE {
        AZPUBSUB_CONFIGURATION_TYPE_GLOBAL,
        AZPUBSUB_CONFIGURATION_TYPE_TOPIC,
        AZPUBSUB_CONFIGURATION_TYPE_SIMPLE
    } AZPUBSUB_CONFIGURATION_TYPE, *PAZPUBSUB_CONFIGURATION_TYPE;

    /**
    An enum defininf the different consumer flags.
    */
    typedef enum _AZPUBSUB_CONSUMER_FLAGS {
        AZPUBSUB_CONSUMER_FLAG_NONE = 0,
        AZPUBSUB_CONSUMER_FLAG_COMMIN_ASYNC = 1
    } AZPUBSUB_CONSUMER_FLAGS, *PAZPUBSUB_CONSUMER_FLAGS;

    /**
    An enum defining the different purge flags.
    */
    typedef enum _AZPUBSUB_PURGE_FLAGS {
        AZPUBSUB_PRODUCER_PURGE_QUEUE = 0x1,
        AZPUBSUB_PRODUCER_PURGE_INFLIGHT = 0x2,
        AZPUBSUB_PRODUCER_PURGE_NONBLOCKING = 0x4
    } AZPUBSUB_PURGE_FLAGS, *PAZPUBSUB_PURGE_FLAGS;

    /**
    * Enum defining global configuration templates
    */
    typedef enum _AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES {
        AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_BALANCED = 0,
        AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_LOW_LATENCY,
        AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_HIGH_THROUGHPUT
    } AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES, *PAZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES;

    typedef struct _AZPUBSUB_BROKER_DATA_ITEM {
        DWORD cbSize;
        DWORD Id;
        DWORD Port;
        DWORD cbHost;
        LPCSTR Host;
    } AZPUBSUB_BROKER_DATA_ITEM, *PAZPUBSUB_BROKER_DATA_ITEM;

    typedef struct _AZPUBSUB_BROKER_DATA {
        DWORD cbSize;
        DWORD BrokerId;
        DWORD cbBrokerName;
        LPCSTR BrokerName;
        DWORD BrokerCount;
        PAZPUBSUB_BROKER_DATA_ITEM Brokers;
    } AZPUBSUB_BROKER_DATA, *PAZPUBSUB_BROKER_DATA;

    typedef struct _AZPUBSUB_PARTITION_DATA_ITEM {
        DWORD Count;
        PDWORD Items;
    } AZPUBSUB_PARTITION_DATA_ITEM, *PAZPUBSUB_PARTITION_DATA_ITEM;

    typedef struct _AZPUBSUB_PARTITION_DATA {
        DWORD cbSize;
        DWORD Id;
        DWORD Leader;
        AZPUBSUB_PARTITION_DATA_ITEM Replicas;
        AZPUBSUB_PARTITION_DATA_ITEM Isrs;
        DWORD Error;
    } AZPUBSUB_PARTITION_DATA, *PAZPUBSUB_PARTITION_DATA;

    typedef struct _AZPUBSUB_TOPIC_DATA
    {
        DWORD cbSize;
        DWORD cbTopic;
        LPCSTR Topic;
        DWORD ErrorCode;
        DWORD PartitionCount;
        PAZPUBSUB_PARTITION_DATA Partitions;
    } AZPUBSUB_TOPIC_DATA, *PAZPUBSUB_TOPIC_DATA;

    typedef struct _AZPUBSUB_TOPIC_DATA_LIST {
        DWORD Count;
        PAZPUBSUB_TOPIC_DATA Data;
    } AZPUBSUB_TOPIC_DATA_LIST, *PAZPUBSUB_TOPIC_DATA_LIST;

    typedef void
    (WINAPI* AZPUBSUB_CLIENT_TOKEN_CALLBACK) ();

    typedef DWORD
    (WINAPI* AZPUBSUB_AUTH_CALLBACK) (
        LPCSTR lpszBrokerName,
        DWORD brokerId,
        LPCSTR lpszAcls,
        INT depth,
        const BYTE* buffer,
        size_t size);

    typedef DWORD
    (WINAPI* AZPUBSUB_CERT_CALLBACK)(
        BYTE** context,
        uint32_t* size,
        PVOID cbContext);

    typedef struct _AZPUBSUB_MESSAGE_HEADER
    {
        LPCSTR pszKey;
        PBYTE pValue;
        DWORD cbValue;
    } AZPUBSUB_MESSAGE_HEADER, *PAZPUBSUB_MESSAGE_HEADER;

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubErrorToString(
        RdKafka::ErrorCode code,
        LPSTR* error,
        PDWORD cchError);

    AZPUBSUB_API
    VOID
    WINAPI
    AzPubSubSetNativeMetricLogging(
        DWORD isMetricLogging,
        DWORD isHotPathMetricLogging);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetNativeMetricLogging();

    AZPUBSUB_API
    HCLIENT 
    WINAPI
    AzPubSubClientInitialize(
        AZPUBSUB_LOG_CALLBACK logRoutine,
        PVOID context);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubClientClose(
        HCLIENT hClient);

    AZPUBSUB_API
    HCONFIG
    WINAPI
    AzPubSubCreateConfiguration(
        HCLIENT hClient,
        AZPUBSUB_CONFIGURATION_TYPE type,
        AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES globalConfigTemplate);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSetSecurity(
        HCONFIG hConfig,
        AZPUBSUB_SECURITY_TYPE type,
        DWORD flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetSecurityFlags(
        HCONFIG hConfig, 
        DWORD* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetSecurityType(
        HCONFIG hConfig, 
        PAZPUBSUB_SECURITY_TYPE pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSetConnection(
        HCONFIG hConfig,
        LPCSTR environment,
        LPCSTR cluster,
        LPCSTR endpoint,
        LPCSTR dnsServerIp,
        AZPUBSUB_SECURITY_TYPE type,
        DWORD flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubAddStringConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        LPCSTR lpszValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubAddBooleanConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        bool bValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubAddIntegerConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        int iValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubAddBinaryConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        const BYTE* value,
        DWORD size);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetStringConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        LPSTR* lpszValue,
        DWORD cchSize,
        PDWORD cchRequiredSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetIntConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        int* piValue);
	
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetBoolConfiguration(
        HCONFIG hConfig,
        LPCSTR lpcszName,
        BOOL* piValue);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSetAuthCallBack(
        HCONFIG hConfig,
        AZPUBSUB_AUTH_CALLBACK sslCertificateVerifyCallback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSetCertCallBack(
        HCONFIG hConfig,
        AZPUBSUB_CERT_CALLBACK certCallback,
        PVOID cbContext);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubFreeConfiguration(
        HCONFIG hConfig);

    AZPUBSUB_API
    HPRODUCER
    WINAPI
    AzPubSubOpenProducer(
        HCONFIG hConfig,
        AZPUBSUB_EVENT_CALLBACK eventCallback,
        AZPUBSUB_METRICS_CALLBACK metricsCallback);

    AZPUBSUB_API
    HPRODUCER
    WINAPI
    AzPubSubOpenSimpleProducer(
        HCONFIG hConfig,
        AZPUBSUB_SECURITY_TYPE type,
        LPCSTR environment,
        LPCSTR cluster,
        LPCSTR endpoint);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerConnect(HPRODUCER hProducer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSendMessage(
        HPRODUCER hProducer,
        HPRODUCERTOPIC hTopic,
        const BYTE* key,
        DWORD keyLen,
        const BYTE* payload,
        size_t payloadLen,
        int32_t partition,
        int64_t timestamp,
        PVOID pHeaders,
        AZPUBSUB_MESSAGE_CALLBACK callback,
        PVOID context);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSendMessageEx(
        HPRODUCER hProducer,
        LPCSTR topic,
        LPCSTR key,
        int* partition,
        const BYTE* payload,
        size_t payloadLen,
        HRESPONSE* pResponse);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubSendBatchMessage(
        HPRODUCER hProducer,
        LPCSTR topic,
        HBATCH hBatch,
        HRESPONSE* pResponse);

    AZPUBSUB_API
    int
    WINAPI
    AzPubSubProducerPoll(
        HPRODUCER hProducer,
        ULONGLONG timeout);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerFlush(
        HPRODUCER hProducer,
        ULONGLONG timeout);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerPurge(
        HPRODUCER hProducer,
        AZPUBSUB_PURGE_FLAGS flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerGetMetadata(
        HPRODUCER hProducer,
        ULONGLONG timeout,
        LPCSTR lpszTopic,
        HMETADATA* pMetadata);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerRegisterTokenRefresh(
        HPRODUCER hProducer,
        AZPUBSUB_TOKEN_REFRESH_CALLBACK callback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerGetConnected(
        HPRODUCER hProducer,
        PBYTE connectedStatus);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerClose(
        HPRODUCER hProducer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerSetGetTokenCallback(
        HPRODUCER hProducer,
        AZPUBSUB_CLIENT_TOKEN_CALLBACK getTokenCallback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerSetTokenInfo(
        HPRODUCER hProducer,
        LONGLONG leftTimeSinceUnixEpochInMs,
        LPCSTR id,
        LPCSTR base64Token);

    AZPUBSUB_API
    HCONSUMER
    WINAPI
    AzPubSubOpenConsumer(
        HCONFIG hConfig,
        LPCSTR groupId,
        AZPUBSUB_EVENT_CALLBACK callback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerConnect(HCONSUMER hConsumer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerRegisterRebalance(
        HCONSUMER hConsumer,
        AZPUBSUB_REBALANCE_CALLBACK callback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerUnregisterRebalance(
        HCONSUMER hConsumer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerRegisterCommit(
        HCONSUMER hConsumer,
        AZPUBSUB_COMMIT_CALLBACK callback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerUnregisterCommit(
        HCONSUMER hConsumer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerSubscribe(
        HCONSUMER hConsumer,
        LPCSTR* ppTopics,
        size_t nTopics);

    AZPUBSUB_API 
    DWORD
    WINAPI
    AzPubSubConsumerUnsubscribe(
        HCONSUMER hConsumer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerAssign(
        HCONSUMER hConsumer,
        HTOPICPARTITIONLIST hList);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerUnassign(
        HCONSUMER hConsumer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerConsume(
        HCONSUMER hConsumer,
        ULONGLONG timeout,
        HMESSAGE* pHmessage);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerCommitAssignment(
        HCONSUMER hConsumer,
        AZPUBSUB_CONSUMER_FLAGS flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerCommit(
        HCONSUMER hConsumer,
        HMESSAGE hMessage,
        AZPUBSUB_CONSUMER_FLAGS flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerCommitEx(
        HCONSUMER hConsumer,
        HTOPICPARTITIONLIST hList,
        AZPUBSUB_CONSUMER_FLAGS flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerModifyDebugState(
        HCONSUMER hConsumer,
        INT logLevel,
        LPCSTR debugKey);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubProducerModifyDebugState(
        HPRODUCER hProducer,
        INT logLevel,
        LPCSTR debugKey);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerCommitted(
        HCONSUMER hConsumer,
        HTOPICPARTITIONLIST hTpList,
        ULONGLONG timeout);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerQueryWatermarkOffsets(
        HCONSUMER hConsumer,
        HTOPICPARTITION hTopicPartition,
        ULONGLONG timeout,
        PULONGLONG pLow,
        PULONGLONG pHigh);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerSeek(
        HCONSUMER hConsumer,
        HTOPICPARTITION hTopicPartition,
        ULONGLONG timeout);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerGetMetadata(
        HCONSUMER hConsumer,
        ULONGLONG timeout,
        LPCSTR lpszTopic,
        HMETADATA* pMetadata);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerRegisterTokenRefresh(
        HCONSUMER hConsumer,
        AZPUBSUB_TOKEN_REFRESH_CALLBACK callback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerGetConnected(
        HCONSUMER hConsumer,
        PBYTE connectedStatus);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerClose(
        HCONSUMER hConsumer);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerSetGetTokenCallback(
        HCONSUMER hConsumer,
        AZPUBSUB_CLIENT_TOKEN_CALLBACK getTokenCallback);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerSetTokenInfo(
        HCONSUMER hConsumer,
        LONGLONG LeftTimeSinceUnixEpochInMs,
        LPCSTR Id,
        LPCSTR Base64Token);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetErrorMessage(
        HMESSAGE hMessage,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetErrorCode(
        HMESSAGE hMessage,
        int32_t* pErrorCode);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetTopicName(
        HMESSAGE hMessage,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetPartition(
        HMESSAGE hMessage,
        int32_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetPayload(
        HMESSAGE hMessage,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetKey(
        HMESSAGE hMessage,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetKeyAsString(
        HMESSAGE hMessage,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetOffset(
        HMESSAGE hMessage,
        int64_t* pValue);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetTimeStamp(
        HMESSAGE hMessage,
        RdKafka::MessageTimestamp* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetLatency(
        HMESSAGE hMessage,
        int64_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetStatus(HMESSAGE hMessage, int32_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageGetHeaders(HMESSAGE hMessage, PVOID* ppHeader);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubMessageClose(
        HMESSAGE hMessage);

    AZPUBSUB_API
    HCONNECTION
    WINAPI
    AzPubSubOpenConnection(
        LPCWSTR lpszEnvironment,
        LPCWSTR lpszCluster,
        AZPUBSUB_SECURITY_TYPE type,
        AZPUBSUB_CONNECTION_FLAGS flags);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubCloseConnection(
        HCONNECTION hConnection);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetConnection(
        HCONNECTION hConnection,
        LPWSTR* lpszConnection,
        LPDWORD pcchConnection);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubResponseGetStatusCode(
        HRESPONSE hResponse,
        int32_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubResponseGetMessage(
        HRESPONSE hResponse,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubResponseGetSubStatusCode(
        HRESPONSE hResponse,
        int32_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubResponseClose(
        HRESPONSE hResponse);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetBrokerMetadata(
        HMETADATA hMetadata,
        PAZPUBSUB_BROKER_DATA* ppList,
        PDWORD pcbList);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetTopicMetadata(
        HMETADATA hMetadata,
        PAZPUBSUB_TOPIC_DATA_LIST* ppList,
        PDWORD pcbList);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubConsumerOffsetStore(
        HCONSUMER hConsumer, 
        HTOPICPARTITIONLIST hList);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubCloseMetadata(
        HMETADATA hMetadata);

    AZPUBSUB_API
    HPRODUCERTOPIC
    WINAPI
    AzPubSubOpenProducerTopic(
        HPRODUCER hProducer,
        LPCSTR lpszTopic,
        HCONFIG hTConfig,
        AZPUBSUB_TOPIC_PARTITIONER partitioner);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubPProducerTopicPartitionAvailable(
        HPRODUCERTOPIC hProducerTopic,
        DWORD partition,
        PBYTE pAvailable);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubPProducerTopicStoreOffset(
        HPRODUCERTOPIC hProducerTopic,
        DWORD partition,
        ULONGLONG offset);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubCloseProducerTopic(
        HPRODUCERTOPIC hTopic);

    AZPUBSUB_API
    HTOPICPARTITION
    WINAPI
    AzPubSubCreateTopicPartition(
        LPCSTR lpszTopic,
        LONG partition);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubTopicPartitionGetTopicName(
        HTOPICPARTITION hTopicPartition,
        BYTE** ppBuffer,
        DWORD size,
        PDWORD pReqSize);
    
    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubTopicPartitionGetPartition(
        HTOPICPARTITION hTopicPartition,
        int32_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubTopicPartitionGetOffset(
        HTOPICPARTITION hTopicPartition,
        int64_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubTopicPartitionGetErrorCode(
        HTOPICPARTITION hTopicPartition,
        int32_t* pValue);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubTopicPartitionSetOffset(
        HTOPICPARTITION hTopicPartition,
        LONGLONG offset);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubCloseTopicPartition(
        HTOPICPARTITION hTopicPartition);

    AZPUBSUB_API
    HTOPICPARTITIONLIST
    WINAPI
    AzPubSubCreateTopicPartitionList();

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetTopicPartitionListSize(
        HTOPICPARTITIONLIST hList,
        PDWORD pSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubGetTopicPartitionFormList(
        HTOPICPARTITIONLIST hList,
        DWORD index,
        HTOPICPARTITION* pTopicPartition);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubAddTopicPartitionToList(
        HTOPICPARTITIONLIST hList,
        HTOPICPARTITION hTopicPartition);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubCloseTopicPartitionList(
        HTOPICPARTITIONLIST hList);

    AZPUBSUB_API
    HBATCH
    WINAPI
    AzPubSubCreateBatch();

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubBatchAddMessage(
        HBATCH hBatch,
        PCSTR key,
        const BYTE* payload,
        DWORD cbPayload);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubBatchGetBuffer(
        HBATCH hBatch,
        BYTE** ppbuffer,
        DWORD size,
        PDWORD pReqSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubBatchGetSize(
        HBATCH hBatch,
        PDWORD pSize);

    AZPUBSUB_API
    DWORD
    WINAPI
    AzPubSubCloseBatch(
        HBATCH hBatch);

    AZPUBSUB_API
    INT
    WINAPI
    AzPubSubGetVarUInt32Length(
        UINT value);

    AZPUBSUB_API
    INT
    WINAPI
    AzPubSubEncodeVarUInt32(
        BYTE* pData, 
        UINT size, 
        UINT value, 
        UINT index);

    AZPUBSUB_API
    INT
    WINAPI
    AzPubSubGetVarUInt64Length(
        ULONGLONG value);

    AZPUBSUB_API
    INT
    WINAPI
    AzPubSubEncodeVarUInt64(
        BYTE* pData, 
        UINT size, 
        ULONGLONG value, 
        UINT index);

    AZPUBSUB_API
    UINT
    WINAPI
    AzPubSubDecodeVarUInt32(
        const BYTE* pData, 
        UINT size, 
        UINT* pIndex);

    AZPUBSUB_API
    ULONGLONG
    WINAPI
    AzPubSubDecodeVarUInt64(
        const BYTE* pData, 
        UINT size, 
        UINT* pIndex);

    AZPUBSUB_API
    UINT
    WINAPI
    AzPubSubEncodeZigzag32(
        INT value);

    AZPUBSUB_API
    ULONGLONG
    WINAPI
    AzPubSubEncodeZigzag64(
        LONGLONG value);

    AZPUBSUB_API
    INT
    WINAPI
    AzPubSubDecodeZigzag32(
        UINT value);

    AZPUBSUB_API
    LONGLONG
    WINAPI
    AzPubSubDecodeZigzag64(
        ULONGLONG value);

#ifdef __cplusplus
} // extern "C"
#endif

#endif //__cplusplus
#endif //_MSC_VER
#endif

#endif // _AZ_PUBSUB_
