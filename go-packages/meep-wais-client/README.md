# Go API client for client

WLAN Access Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC028 WAI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.01.01_60/gs_MEC028v020101p.pdf) <p>[Copyright (c) ETSI 2020](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-wais](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-wais) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about WLAN access information in the network <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_

## Overview
This API client was generated by the [swagger-codegen](https://github.com/swagger-api/swagger-codegen) project.  By using the [swagger-spec](https://github.com/swagger-api/swagger-spec) from a remote server, you can easily generate an API client.

- API version: 2.1.1
- Package version: 1.0.0
- Build package: io.swagger.codegen.languages.GoClientCodegen

## Installation
Put the package under your project folder and add the following in import:
```golang
import "./client"
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost/wai/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DefaultApi* | [**ApInfoGET**](docs/DefaultApi.md#apinfoget) | **Get** /queries/ap/ap_information | 
*DefaultApi* | [**StaInfoGET**](docs/DefaultApi.md#stainfoget) | **Get** /queries/sta/sta_information | 
*DefaultApi* | [**SubscriptionLinkListSubscriptionsGET**](docs/DefaultApi.md#subscriptionlinklistsubscriptionsget) | **Get** /subscriptions/ | 
*DefaultApi* | [**SubscriptionsDELETE**](docs/DefaultApi.md#subscriptionsdelete) | **Delete** /subscriptions/{subscriptionId} | 
*DefaultApi* | [**SubscriptionsGET**](docs/DefaultApi.md#subscriptionsget) | **Get** /subscriptions/{subscriptionId} | 
*DefaultApi* | [**SubscriptionsPOST**](docs/DefaultApi.md#subscriptionspost) | **Post** /subscriptions/ | 
*DefaultApi* | [**SubscriptionsPUT**](docs/DefaultApi.md#subscriptionsput) | **Put** /subscriptions/{subscriptionId} | 


## Documentation For Models

 - [ApAssociated](docs/ApAssociated.md)
 - [ApIdentity](docs/ApIdentity.md)
 - [ApInfo](docs/ApInfo.md)
 - [ApLocation](docs/ApLocation.md)
 - [BeaconReport](docs/BeaconReport.md)
 - [BeaconRequestConfig](docs/BeaconRequestConfig.md)
 - [BssLoad](docs/BssLoad.md)
 - [ChannelLoadConfig](docs/ChannelLoadConfig.md)
 - [CivicLocation](docs/CivicLocation.md)
 - [DmgCapabilities](docs/DmgCapabilities.md)
 - [EdmgCapabilities](docs/EdmgCapabilities.md)
 - [ExtBssLoad](docs/ExtBssLoad.md)
 - [GeoLocation](docs/GeoLocation.md)
 - [HeCapabilities](docs/HeCapabilities.md)
 - [HtCapabilities](docs/HtCapabilities.md)
 - [InlineResponse200](docs/InlineResponse200.md)
 - [InlineResponse2001](docs/InlineResponse2001.md)
 - [InlineResponse2002](docs/InlineResponse2002.md)
 - [InlineResponse2003](docs/InlineResponse2003.md)
 - [InlineResponse201](docs/InlineResponse201.md)
 - [InlineResponse400](docs/InlineResponse400.md)
 - [InlineResponse403](docs/InlineResponse403.md)
 - [Link](docs/Link.md)
 - [NeighborReport](docs/NeighborReport.md)
 - [OptionalSubelement](docs/OptionalSubelement.md)
 - [ProblemDetails](docs/ProblemDetails.md)
 - [StaDataRate](docs/StaDataRate.md)
 - [StaIdentity](docs/StaIdentity.md)
 - [StaInfo](docs/StaInfo.md)
 - [StaStatistics](docs/StaStatistics.md)
 - [StaStatisticsConfig](docs/StaStatisticsConfig.md)
 - [StatisticsGroupData](docs/StatisticsGroupData.md)
 - [Subscription](docs/Subscription.md)
 - [Subscription1](docs/Subscription1.md)
 - [SubscriptionLinkList](docs/SubscriptionLinkList.md)
 - [SubscriptionPost](docs/SubscriptionPost.md)
 - [SubscriptionPost1](docs/SubscriptionPost1.md)
 - [TimeStamp](docs/TimeStamp.md)
 - [VhtCapabilities](docs/VhtCapabilities.md)
 - [WanMetrics](docs/WanMetrics.md)
 - [WlanCapabilities](docs/WlanCapabilities.md)


## Documentation For Authorization

## OauthSecurity
- **Type**: OAuth
- **Flow**: application
- **Authorization URL**: 
- **Scopes**: 
 - **all**: Single oauth2 scope for API

Example
```golang
auth := context.WithValue(context.Background(), sw.ContextAccessToken, "ACCESSTOKENSTRING")
r, err := client.Service.Operation(auth, args)
```

Or via OAuth2 module to automatically refresh tokens and perform user authentication.
```golang
import "golang.org/x/oauth2"

/* Perform OAuth2 round trip request and obtain a token */

tokenSource := oauth2cfg.TokenSource(createContext(httpClient), &token)
auth := context.WithValue(oauth2.NoContext, sw.ContextOAuth2, tokenSource)
r, err := client.Service.Operation(auth, args)
```

## Author


