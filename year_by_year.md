j# Current Year_by_year

One of the nicest graphs to look at is going to be of your revenue growing, One of the best ways (and eaiest to implemenet) is going to be a a yearly compairison. This is exactly what Year_by_Year is. grabbing all of the orders by week and summing their revenue and putting that on a chart.

I need to write down what I'm thinking about the best implement the year by year so I can improve my design.


## Previous implementations

I actualy have a couple implementations of this, all in different languages... First I implemented this in python with Flaskas a Proof of Concept running on my computer. This was the worst solution because it I actually did the caclulation every time. The calculation was grabbing every record in the last 4 years summing the weeks and then sending over the data to the client. FYI, I dont wanna talk too much about the client because front end is anoyying. 

I think its odvious that that isn't a great idea, especially in python. When you have thousands of records a year then doing additional arithmetic on the data to send it out. It was very slow, especially once I put it onto a raspberry pi as the server. 

I was also having compatatbility issues with the pi due to its age and so I though you know what? Lets get rusty and rewrite it in rust. Now im a smart kid, Ive learnt a little bit of rust and its very satisfying when the problem is pretty easy. Now my first sollution was pretty interesting and it can be fairly fast! thats what I needed for the pi. Basically I implemented a cached solution. So I had the object that I created, basically all the data that I was going to send off cached in memory. Ill be honest I don't even know how I did because Dealing with mutexes and Arc's was rough. I didnt have a very good understanding of that when I started. I learned a lot but there is still so much to learn about rusts, parrellelisam. The think is the data wasn't changing very much, especially the old data. Really the only thing that was changing was the newest weeks. WIth my implementation I could easly make it just update the cache with the new data but I had it update everything. either way. I dont know if this is the write solution.

## GO lang

Now I'm a new man, Ive started learning golang. the python of compiled languages. Now I'm not very good at it and I'm still learning a lot but its also a straight forward lang as well. I'm much quicker at just implementing somthing quick and able to improve on it much quicker then when I was battleing the rust compiler. (Probably skill issues I know) basically How I ended up here is I left the backend server running and then I also created another backend server in go with gin to run the new idead that I had. Dispachers! I got it implemented pretty quickly and there are some good things about it but there are also a few things that I need to change but I need to go over it with one of my bossmen. Anyways this isn't really a Dispachers file but talking about Year_by_year.

### Implementation

Thinking bout it I need to go even more indepth with the server config. The rust server would first have an empty cache witch is the first thing it would look for. Then it would run the expensive calcuation to figure out what the data was, Then it would have a timestamp of when that was done, It would then return that data that would be saved in memory anytime it was called. but first it would check the database once for a new record in the "db updated" table that would just save a timestamp if I updated the db. It was pretty quick when it came down to it but I dont think this is the best method for a webserver with huge requirements and an everchanging db.

### New Implementation

I'm think the new implementation Is going to use another new table! But I think that it makes way more sence too. I think it makes the most sence to have all of those numbers that I caculate just stored into a db as | ID | Revenue | Week | Quarter | Year | I think that is all the information that I need to send because then all I need to do is select * all the records and convert them into my data structure that I'll be using to send the json but also ill just add like a clause to check what week it is currently, and any week or weeks that are not added we can search that from the table of records and sum up that week and get that information we need. I wanna add some psudocode here just to make my point come accross

```go
// models

type YearlyRevenue struct {

    Week         String
    2021_Revenue float
    2022_Revenue float
    2023_Revenue float
    2024_Revenue float

}

```

```go
//handler

db_connection, err = make_connection_to_db()

data, err := db_connection.query(SELECT * FROM year_by_year_rev Groupby Week) // should return 3 or 4 records for each week

data, err := getLatestData(db_connection, data) // function responsible for getting the newest data 
// isnt added to the other db yet and adding that info to the dataset.

c.JSON(200, data)
```
SO i think that I need to create the db before I can finish the server access that will make it a little better, so I wrote a quick program to access all the records and turn them into a Records slice and then add them to a new table before.

Okay so the first kinda implementation that I did it seems that we can't access two tables at the same time when we have rows open and so we're going to create a slice

## Another Try

After trying to get this going and working It seems that I needed to just utilize a map for the data so I can get all of the old data and then we arrive at a week that hasn't happend yet... we just don't add it to the dictionary. 

So far I got the code working and it is returning the proper information, I still havent implemented the logistics side. (This should be easy because its really the same exact code but just a different database table). I think I would like to create a drop down table that will be like radio options where you can select which years that you would like to compair. I think I'm just going to continue to send all of the data and then add a drop end to "filter" the data on the front end side.
