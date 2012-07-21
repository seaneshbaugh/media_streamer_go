# Media Streamer (Go Implementation)

A simple music server using written in Go. This is more or less a direct port of the [Ruby version](https://github.com/seaneshbaugh/media_streamer) of this application. I suspect it hardly counts as idiomatic Go. I also suspect that it currently doesn't fully utilize Go's capabilities. While the Ruby version was made to scratch an itch, this implementation is just for fun.

Go Setup
--------

Here's some instructions to get you started with Go. This is my first attempt at anything more serious than a "Hello World" program in Go. As such I ran into some confusion getting this project set up correctly. If you've already got a working Go environment you can safely ignore this. If you're like me and new to the language and its conventions you probably should read this. There's [instructions online here](http://golang.org/doc/code.html)  but I found them to be not only unclear, they were in a format such that I had to tease out the relevant details (i.e., not a list of instructions). Here's how I eventually got everything set up:

### Windows 7

Some specifics may be different on older or newer versions of Windows.

* Install Go from the [.msi installer](http://code.google.com/p/go/downloads/list?q=OpSys-Windows+Type%3DInstaller). You can download it as a Zip archive, but why would you do that?
* Create a directory for your projects if you haven't already. I keep all my projects in a folder directory *Projects* in my *Documents* directory. In my case it's *C:\Users\Sean\Documents\Projects*.
* Create a directory specifically for your Go projects. In my case it's *C:\Users\Sean\Documents\Projects\go*.
* Create the following directories in your Go projects directory: *bin*, *pkg*, and *src*. *src* is where our Go project source files will reside.
* Add the GOPATH environment variable to your user environment variables.
    * In the start menu right click the "Computer" menu item and then click "Properties".
    * Click "Advanced system settings" on the left (it should be the last option in the list).
    * Click "Environment Variables".
    * Under "User variables for *username*" click "New...".
    * Make the variable name be "GOPATH".
    * Make the variable value be the path to your Go projects directory. In my case it's *C:\Users\Sean\Documents\Projects\go*.
    * Keep clicking "OK" until you can't click it any more.
    * If you have any open command prompts you will need to close them and re-open them for the changes to take effect.
    * In the command prompt test it out with `echo %GOPATH%`.
* cd into your %GOPATH%/src directory and create a folder for your project(s).
* Create directories in your project directory for each package you are making for your project.
* Add the source files for each package.

To be clear: the second and third steps are mostly just my opinion.

### OSX Lion

Some specifics may be different depending on your environment.

* Install Go from the [.pkg installer](http://code.google.com/p/go/downloads/list?q=OpSys-Darwin).
* Create a directory for your Go projects if you haven't already. I keep mine *~/go*.
* Create the following directories in your Go projects directory: *bin*, *pkg*, and *src*. *src* is where our Go project source files will reside.
* Add the GOPATH environment variable to your user environment variables. This will depend on your shell. Add the following to your .bashrc file (or if you're using ZSH your .zshrc file): GOPATH=$HOME/go
* Source the changes with source .bashrc (or source .zshrc) in the terminal.
* Make sure it's working with `echo $GOPATH`.
* cd into your $GOPATH/src directory directory and create a folder for your project(s).
* Create directories in your project directory for each package you are making for your project.
* Add the source files for each package.

### Linux

More on this later... but it's probably pretty close to OSX.

Installation
------------

		$ cd $GOPATH/src (%GOPATH%\src on Windows)
    $ git clone git@github.com:seaneshbaugh/media_streamer_go.git media_streamer_go
    $ cd media_streamer_go/webserver
    $ go build webserver.go

Usage
-----

    $ $GOPATH/src/media_streamer_go/webserver/webserver (%GOPATH%\src\media_streamer_go\webserver\webserver.exe on Windows)

Options
-------

* -d, Sets the "public" directory for the server. This is the directory from which static assets will be delivered. The default is *./public*.
* -p, Sets the port the server will listen on. The default is 4568.
* -c, Enables or disables GZip Compression. Must be used like this: -c=0, -c=1, -c=true, -c=false, -c=TRUE, or -c=FALSE. The default is true.
* -m, Sets the "media" directory for the server. This is the directory where the server will look for music files. The default is */Users/seshbaugh/Music/iTunes/iTunes Media/*. You will definitely need to change this before building.

Notes
-----

Many thanks to [Alexis Robert](https://github.com/alexisrobert) whose [simple Go webserver](https://gist.github.com/982674) gave me a working example to start from.