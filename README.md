# Surfs [![Actions Status](https://github.com/maybetheresloop/surfs/workflows/Go/badge.svg)](https://github.com/maybetheresloop/surfs)

A simple distributed file store in Go, based on UCSD CSE 124's SurfStore assignment,
originally implemented in Java.

## Project Status
This project, like the Keychain project that it uses for metadata persistence, is being
written for learning purposes. As with Keychain, the goal is to implement basic functionality
first and then move onto more advanced features.

## Internals

The following is an overview of how Surfs works internally, taken from the
SurfStore CSE 124 assignment description:

> A file in SurfStore is broken into an ordered sequence of one or more blocks. Each 
> block is of uniform size (4KB), except for the last block in the file, which may be
> smaller than 4KB (but must be at least 1 byte large).
>
> For each block, a hash value is generated using the SHA-256 hash function. This set of 
> hash values, in order, represents the file, and is referred to as the hashlist. Note that 
> if you are given a block, you can compute its hash by applying the SHA-256 hash function 
> to the block. This also means that if you change data in a block the hash value will change 
> as a result. To update a file, you change a subset of the bytes in the file, and recompute 
> the hashlist. Depending on the modification, at least one, but perhaps all, of the hash values 
> in the hashlist will change.