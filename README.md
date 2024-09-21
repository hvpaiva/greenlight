# Greenlight API Project

This repository contains a fully implemented movie catalog API project following the steps laid out in the book *Let's Go Further* by Alex Edwards. The project demonstrates advanced patterns in API development using Go, and covers topics such as authentication, authorization, database handling, and more.

While the code here follows the content of the book, I have made some personal adjustments to structure, error handling, and naming conventions to better fit my style and preferences. Nevertheless, all the core features from the book are fully implemented.

If you're considering purchasing *Let's Go Further*, I highly recommend it. The book provides invaluable theoretical insights and deeper explanations that complement the practical code here.

## Project Overview

The **Greenlight API** provides a full-featured backend for managing a movie catalog with the following features:

- **Authentication**: Secure login, JWT tokens, and permission-based authorization.
- **Movie Management**: Full CRUD operations for movie entries.
- **Filtering & Pagination**: Support for filtering, sorting, and paginating API responses.
- **Metrics**: Exposes application metrics, such as memory usage and request counts.
- **Security**: Implements rate limiting and CORS management for enhanced API security.

## Key Features

This project covers the following advanced topics:

- **JSON Request Handling**: Parsing and validating incoming JSON requests.
- **Database Management**: SQL migrations, connection pooling, and query optimizations.
- **Authentication & Authorization**: Supports both JWTs and stateful sessions, along with granular permission checks.
- **Background Processing**: Safely manages background tasks with goroutines.
- **Email Services**: Sends user activation and password reset emails.
- **Rate Limiting**: Protects API endpoints by limiting the number of requests per client.
- **CORS**: Controls cross-origin requests for better API security.

## Acknowledgements

This project is based on the book *Let's Go Further* by Alex Edwards. While I've implemented the entire project described in the book, the theoretical discussions and detailed explanations in the book provide valuable context and understanding beyond the code itself. If you're looking to deepen your Go knowledge, the book is a must-read.

You can learn more about the book [here](https://lets-go-further.alexedwards.net).
