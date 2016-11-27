# Golang URL Minifier v0.0.2

[Check it out at http://r8r.org](http://r8r.org)

A URL minifier written in Go.

  - Enter long URL
  - Receive small URL
  - Magic?

You can also:
  - Use the API to integrate your apps

<img src="https://raw.githubusercontent.com/nickvellios/Golang-URL-Minifier/master/ss1.png" alt="First Screenshot" width="45%" height="45%">
<img src="https://raw.githubusercontent.com/nickvellios/Golang-URL-Minifier/master/ss2.png" alt="Second Screenshot" width="45%" height="45%">

### Tech

Golang URL Minifier relies on open source projects to work properly:

* [pq] - A pure Go postgres driver for Go's database/sql package
* [Twitter Bootstrap] - For quick UI creation
* [jQuery] - Because Stackoverflow said I had to?

And of course Golang URL Minifier itself is open source with a [public repository][gomin]
 on GitHub.

### Installation

Golang URL Minifier requires [pq](https://github.com/lib/pq) to run.

```sh
$ go get github.com/lib/pq
```

### API Usage

Limited to 10 requests per hour per IP address

Example Request (via POST only):

```sh
curl --request POST 'http://r8r.org/generate/' --data "url=http://www.golang.org"
```

Example JSON Response:

```sh
{
	"url": "http://r8r.org/dvHhd"
	"error": ""
}
```

### Todos

 - Admin dashboard
 - User accounts to remove throttle
 - Traffic analytics (Update:  [stats] now exists, lots of work to do here!)
 - ???
 - Profit!

License
----

None.  For more information, please refer to <http://unlicense.org/>


**Free Software, Hell Yeah!**

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)


   [pq]: <https://github.com/lib/pq>
   [Twitter Bootstrap]: <http://getbootstrap.com/>
   [jQuery]: <https://jquery.com/>
   [gomin]: <https://github.com/nickvellios/Golang-URL-Minifier>
   [stats]: <http://r8r.org/stats/>
