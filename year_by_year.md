# Current Year_by_year

One of the nicest graphs to look at is going to be of your revenue growing, One of the best ways (and eaiest to implemenet) is going to be a a yearly compairison. This is exactly what Year_by_Year is. grabbing all of the orders by week and summing their revenue and putting that on a chart.

I need to write down what I'm thinking about the best implement the year by year so I can improve my design.


## Previous implementations

I actualy have a couple implementations of this, all in different languages... First I implemented this in python with Flask as a Proof of Concept running on my computer. This was the worst solution because it I actually did the caclulation every time. The calculation was grabbing every record in the last 4 years summing the weeks and then sending over the data to the client. I think its odvious that that isn't a great idea, especially in python. When you have thousands of records a year then doing additional arithmetic on the data to send it out. It was very slow, especially once I put it onto a raspberry pi as the server. 

I was also having compatatbility issues with the pi due to its age and so I though you know what? Lets get rusty and rewrite it in rust. Now im a smart kid, Ive learnt a little bit of rust and its very satisfying when the problem is pretty easy. Now my first sollution was pretty interesting and it can be fairly fast! thats what I needed for the pi. Basically I implemented a cached solution. So I had the object that I created, basically all the data that I was going to send off cached in memory. Ill be honest I don't even know how I did because Dealing with mutexes and Arc's was rough. I didnt have a very good understanding of that when I started. I learned a lot but there is still so much to learn about rusts, parrellelisam. The think is the data wasn't changing very much, especially the old data. Really the only thing that was changing was the newest weeks. WIth my implementation I could easly make it just update the cache with the new data but I had it update everything. either way. I dont know if this is the write solution.

## GO lang

Now I'm a new man, Ive started learning golang. the python of compiled languages. Now I'm not very good at it and I'm still learning a lot but its also a straight forward lang as well. I'm much quicker at just implementing somthing quick and able to improve on it much quicker then when I was battleing the rust compiler. (Probably skill issues I know. Async rust is hard.) basically how I ended up here is I left the backend server running and then I also created another backend server in go with gin to run the new idead that I had. Dispachers! I got it implemented pretty quickly and there are some good things about it but there are also a few things that I need to change but I need to go over it with one of my bossmen. Anyways this isn't really a Dispachers file but talking about Year_by_year.

### Implementation

Thinking bout it I need to go even more indepth with the server config. The rust server would first have an empty cache witch is the first thing it would look for. Then it would run the expensive calcuation to figure out what the data was, Then it would have a timestamp of when that was done, It would then return that data that would be saved in memory anytime it was called. but first it would check the database once for a new record in the "db updated" table that would just save a timestamp if I updated the db. It was pretty quick when it came down to it but I dont think this is the best method for a webserver with huge requirements and an everchanging db.

### New Implementation

I'm think the new implementation Is going to use another new table! But I think that it makes way more sence too. I think it makes the most sence to have all of those numbers that I caculate just stored into a db as | ID | Revenue | Week | Quarter | Year | I think that is all the information that I need to send because then all I need to do is select * all the records and convert them into my data structure that I'll be using to send the json but also ill just add like a clause to check what week it is currently, and any week or weeks that are not added we can search that from the table of records and sum up that week and get that information we need. 
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
## Another Try

Now that its all said and done. I've moved off of a rasberry pi and Now I just have an old laptop running all the servers. We have 3 instances running on the same machine. First we have the golang server odviously, but then we have the postgress database running and this is the key to get it all working. firstly I have a direct in memory cache for the quickest retrival. This was pretty easy to implement, but one of the problems came from wanting to get every single record avalable from the last few years. This was not a quick calculation. So what I implemented the postgress database to store each weeks total revenue. I don't want to really create a new table and mess with things on their server so i've pretty much been read only for that one. On my server I've added a table for the year week and revenue. then we search the mcloud for the last database entry too today. that way we are really relaxing the servers and not having to pull tens of thousands of records. Then if we got over a week from their server we create a new database entry to cache. then finally updating our in memory cache so we dont need to do any fetching. So yeah I think this is the best way to do it but i'm just kinda doing what I know and I dont know what I dont know so if there is a better method I would love to hear it. 

Its 2025 now! time flies. One thing I need to figure out is how to edit this data so that it works no matter what year creating that struct worked well for the year the latest year was 2024 but now that its 2025 I need to figure out how to make it so that it works no matter what year it is. I might have to move away from a struct and use a map? I'll have to do some testing.


