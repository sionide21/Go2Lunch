WARNING
=======
This is mega super alpha. It has lots of bugs and is probably very insecure. Use at your own risk.

Go2Lunch
========

Go2Lunch is a tool for organizing lunch groups. The client runs as a command line application with no dependencies. The server stores votes in memory.

## Setup

Go2Lunch ships with a Rakefile to make building easy. If you have rake installed, simply type `rake` to build both the client and server applications.

If you prefer not to use rake, build using the following commands

    6g server.go common.go model.go
    6l -o server server.6
    6g lunch.go common.go update.go
    6l -o lunch lunch.6

## Server

The server requires a config file to operate. The default file is called `config.json` and should be in the same directory as the server. A different file can be specified with the `-c` flag.

### Management

A ruby script called `manage` is included to make server configuration easier. By default, manage will create a config file for you.

#### Adding a user

    ./manage add < keyfile

`keyfile` should be the key provided by a client. It looks like `"name":"LONGSTRINGOFCHARACTERS=="`.

#### Removing a user

    ./manage rm <name>

#### Listing users

    ./manage

## Client

The client must be configured with a name and the server it uses on first run. Use the following command to setup

    lunch -n <yourname> -s <serveraddress>:<port>

This will spit out a large string in the form `"name":"LONGSTRINGOFCHARACTERS=="`. Add this to the server using the manage script.

## Usage

### View current results

    lunch

### Nominate a location

    lunch -a <name of place>

### Vote for a location

    lunch <location #>

### Cancel your vote

    lunch -u

### Volunteer to drive

    lunch -d <available seats>

### Remove a place you nomnated earlier

    lunch -rm <location #>

Note: You can only remove a place if no one has voted for it.
