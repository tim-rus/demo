# Project Overview

This project is a demonstration simulator of an AI-powered ticket routing system for a service center. Its purpose is to generate synthetic customer requests for an imaginary electronics repair center and process them through a realistic distributed pipeline — message broker, RAG-based classifier, and workers — while exposing a live UI for request tracing, load monitoring, failure injection, and token/cost accounting. The simulator is intended to reproduce enterprise-grade operational concerns, including metrics, logging, retries, and idempotency, within a controlled test environment.

# Documentation

The ./README.md is in Russian and for humans only.  

./docs/* are in English and primarily for agents.  

# Git

A new branch for each session.

Commit often, commit small: each commit should be as atomic as possible.

# Rules

use /caveman  
for commit messages use /caveman-commit  
for questions always use widget  

do not scan and do not edit ./prompts  
do not edit ./AGENTS.md
