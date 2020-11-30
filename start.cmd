@echo off

call redis-server.exe

cirrus.exe -config cirrus.toml