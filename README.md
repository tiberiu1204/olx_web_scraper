# Olx web scraper API for finding arable land for sale
### Description
This is a project made for personal use, out of curiosity and interest in the arable land market. The API has multiple filters, and can calculate multiple useful data and relay them to the user, like average and median price points, minimum and maximum price per hectar (RON/ha) offer.
### Usage
Download and install go if you don't already have it. Clone the repository and run the project API:
```
go run
```
The API can be testd with sofware like Postman. An example api call can be found in this repository.
### Filters
There are 6 filters you can apply to your query:

- MinPrice (uint32)  // Minimum price of scraped offers
- MaxPrice (uint32)   // Maximum price of scraped offers
- MinArea (uint32)    // Minimum area of scraped offers
- MaxArea (uint32)    // Maximum area of scraped offers
- MinPPH (float64)    // Minimum price per hectar of scraped offers 
- MaxPPH (float64)    // Maximum price per hectar of scraped offers
