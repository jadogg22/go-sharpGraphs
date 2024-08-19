# Cache !

I'm super stoked because I got access to the actual mcloud db and I can start making my own querys. beforehand I was only able to get the info I needed from parsing csv files. Now I can get the info I need directly from the db. THey are still concerned I was just gonna overwelm the server with a buncha pings so I let them know I could probably have some kind of in memory cache. I can query everything that I need and then store it in a cache. I can then query the cache instead of the db and have the data automatically go stale after like 30 minutes or so. SO if some reason the entier office visited the site at the same time It would only max hit the mcloud serers 1 time and the rest of the time is gonna come from cache.

