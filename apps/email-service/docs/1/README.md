# Basic architecture

## Receiving email

This is the scenario where any one sends an email to anyone on the jlrickert.me
domain. This is the happy path

- User sends email to SES via SMTP protocol (Need to verify this)
- SES responds by triggering a couple of actions. Store the email in S3.
  Triggering an email forwarding lambda to send it off to gmail.com

- Amazon ses

Incoming email pipeline

- user (ex. some one else gmail account) sends email to xyz@jlrickert.me

  SMTP protocol is used to send to SES

- SES forwards email to lambda
- Lambda does a bunch of things

  - Store email in S3
  - Send log to cloud watch
  - ??? Notify SES SMTP thingy for my email

## Sending email

- Gmail sends email to SES via SMTP
- SES forwards email to other user

## SES email configuration

- Email receiving role

  - Role name: basic-email-role
  - Recipient conditions: jlrickert.me
  - Actions: Deliver to s3
    - Bucket: s3-me-jlrickert-email

## SES Forwarding role

- name: _ses-forwarder-role_
- trusted entity: lambda
- inline policy

  - name: ses-forwarder-policy
  - data

    ```json
    {
    	"Version": "2012-10-17",
    	"Statement": [
    		{
    			"Effect": "Allow",
    			"Action": [
    				"logs:CreateLogGroup",
    				"logs:CreateLogStream",
    				"logs:PutLogEvents"
    			],
    			"Resource": "arn:aws:logs:*:*:*"
    		},
    		{
    			"Effect": "Allow",
    			"Action": "ses:SendRawEmail",
    			"Resource": "*"
    		},
    		{
    			"Effect": "Allow",
    			"Action": ["s3:GetObject", "s3:PutObject"],
    			"Resource": "arn:aws:s3:::s3-me-jlrickert-email/*"
    		}
    	]
    }
    ```
