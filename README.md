# Discord Utility Bot

## Configuration
* Edit config.yml
* Create `secrets.yml`:
    ```yaml
    discord:
      token: changeme
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