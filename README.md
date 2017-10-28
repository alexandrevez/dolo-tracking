# dolo-tracking

dolo-tracking is a tool to automate Doloreanne's tracking. This is a tool to fetch contacts (TODO define filter), add a deal to a pipeline, and send a tracking email using Sparkpost template. 

For specifics about the tools, build your copy and run `bin/dolo-tracking --help`.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Installation

### Clone the repo

	mkdir -p $HOME/dev/
	git clone https://github.com/alexandrevez/dolo-tracking.git
	cd dolo-tracking

### macOS

	sudo easy_install pip
	sudo pip install -U pip
	sudo pip install ansible
	ansible-playbook --extra-vars "app_path=`pwd`" -i "localhost," -c local -K ansible/dev/macos/main.yml
	source $HOME/.profile

### Manual installation

See [this gude](docs/INSTALL.md) for manual installation

## Running the app

To build the application, simply run:

	make

You can then launch the applications in `bin` folder depending on which one you want to work with. For example:

	bin/dolo-tracking

Shortcut version:

	make && bin/dolo-tracking


## Running the tests
Running the tests will lint and run the tests for every packages with test coverage. 
	
Coverage report will be located in test/ folder after running the command :

	make test