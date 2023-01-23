# Discord Utility Bot

Add to server: [![Discord](https://img.shields.io/badge/Discord-%235865F2.svg?style=flat&logo=discord&logoColor=white)](https://discord.com/api/oauth2/authorize?client_id=1066509605556977724&permissions=8&scope=applications.commands%20bot)

## Configuration
* Edit config.yml
* Create `secrets.yml`:
    ```yaml
    discord:
      token: Bot changeme
    ```

## Usage
### Commands
#### ublock
Check what resources are blocked from loading on a page by uBlock Origin

### Features
#### Call Summary
* Join channels when they become active
* Stopwatch for total call duration
* Stopwatch for individual speakers
* Post summary to text when everyone leaves or when commanded
#### Multichannel Audio Recording
* Create recording of a voice channel when summoned
* Output to separate audio files per speaker
* Option to output to separate audio outputs live for OBS etc.
    * is there some library for ipc audio?