## Flago
Flago is a feature-flagging solution built on top of redis.Right now flago works only of feature deployments, but could easily be developed to accomodate more flags
Flago is eventually consistent and takes <60s to replicate across data stores.

## Features

1. Don't deploy to any customers (i.e. "inert deploy") 
2. Deploy to one "test" customer 
3. Deploy to a set of "test" customers 
4. Deploy to all alpha customers, where alpha is a standard set of testing customers used by your company.[x]
5. Deploy to all alpha customers + N% of free-tier customers
6. Deploy to all alpha + free-tier customers + N% of premium-tier customers
7.Deploy to everyone 

## Approach
Flago enables complex rollouts by separating the two concerns in feature-flagging:

1. The control plane is where you define your flags, perhaps through a web UI, a CLI tool, an API, or all of the above. Flago lets you define flags with a basic expression syntax that's just powerful enough to cover most use-cases.
2. The data plane is the place in your code that checks whether to use a feature or not. Flago keeps this as simple as possible. You ask Flago whether a flag is enabled for a customer or not. All you get is a true/false, or an error if the flag backend can't be reached.


