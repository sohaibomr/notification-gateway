# Swvl Backend Challenge
In SWVL, we are communicating with our customers via notifications. we are sending
Promo codes to customers via SMS. Also, We are sending to our riders notifications
during their ride like “Your drop-off station is coming”. Customer view all push
notifications in the App. Each customer receives a notification with a preferred
language.

## Specifications:
- Some Notification will be sent by SMS.
- Some Notification will be sent via Push Notification to mobile.
- We have two types of notification
- Group notification which are send as a text notification to a group of users.
- The personalized notification which are send as a specific text notification to each
user.
- The number of requests of providers (SMS, Email) can handle per minute are
limited.

## Requirements:
- Design and implement the notification schema.
- Implement notification service which handles the database changes and sends
notifications for customers.
- Service should be run by docker-compose up.
- Document your APIs.
- Test your code using unit test.

# How to Build/Run:
1. Build the docker images of all services with `make build`.

2. Run all serices with `make run` or `docker-compose up`.

3. Run `make clean` to delete all iamges.

# Solution Overview & Architecture
TODO: write some lines about why chosed this design etc
![title](notification-gateway.jpg)

- notification-server
- group-notifier service
- user-notifier service
- notification-forwarder service

## Functional Description
**Step 1:** POST Request is made to the `notification-service` for group notification `/notification/group` or user notification `/notification/user` to the notifications server, which then parses and validates the request object and pushes notification to the `group` or `user` kafka topic based on the NotificationType `type`.

**Step 2:** `group-notifier service` or `user-notifier service` then consumes messages/notifications from their respective kafka topics, recreates the notification by querying the respective DB to enhancing the user or group message with their meta, tokens, message templates and formating the final message.

**Step 3:** `group-notifier service` or `user-notifier service` then pushes the enhanced notification to the `sms` or `push` kafka topic based on the `sendVia` type.

**Step 4:** `notification-forwarder` then consumes the notifications from the `sms` and `push` topics, and can forward the notifications to the Twillo for SMS or Firebase Cloud Messaging (FCM) or Apple Push Notification service(APNS).

**Note:** For `notification-forwarder` logs all the consumed notifications to console.
# API
`notification-service` exposes two endpoints for group and user notifications.
- `/notification/group`
body:
```
{
    "groupId": "1",
    "sendVia": "sms",
    "message": "Hello Swvl, avail 20 % off",
    "type":"group",
    "category":"promo",
    "tags":["newyYear", "winter", "exams"]
    createdAt:"2020-10-03T15:04:05.0000000Z" // accepts standard RFC 3339 format
}
```
- `/notification/user`
body:
```
{
    "userId": "1",
    "sendVia": "sms",
    "message": "Hello Swvl, avail 20 % off",
    "type":"group",
    "category":"promo",
    "tags":["newyYear", "winter", "exams"]
    createdAt:"2020-10-03T15:04:05.0000000Z" // accepts standard RFC 3339 format
}
```

# Fault Tolerance & High availability
# Scalability & Performance
