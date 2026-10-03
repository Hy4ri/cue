local mp = require('mp')
local utils = require('mp.utils')
local ready = false
local restoring = false
local saved = nil

local file = io.open(settings_path, 'r')
if file then
    saved = utils.parse_json(file:read('*a'))
    file:close()
end

local function persist()
    if not ready or restoring then return end
    local value = {
        aid = mp.get_property_native('aid'),
        sid = mp.get_property_native('sid'),
        audio_delay = mp.get_property_number('audio-delay', 0),
        sub_delay = mp.get_property_number('sub-delay', 0)
    }
    if value.aid == nil or value.sid == nil then return end
    local json = utils.format_json(value)
    if saved and utils.format_json(saved) == json then return end
    local output, err = io.open(settings_path, 'w')
    if not output then
        mp.msg.error('Cue: cannot save playback settings: ' .. tostring(err))
        return
    end
    output:write(json)
    output:close()
    saved = value
    mp.msg.info('Cue: saved audio=' .. tostring(value.aid) .. ' subtitles=' .. tostring(value.sid))
end

mp.register_event('start-file', function() ready = false end)
mp.register_event('file-loaded', function()
    restoring = true
    if saved then
        mp.set_property_native('aid', saved.aid)
        mp.set_property_native('sid', saved.sid)
        mp.set_property_number('audio-delay', saved.audio_delay or 0)
        mp.set_property_number('sub-delay', saved.sub_delay or 0)
        mp.msg.info('Cue: restored audio=' .. tostring(saved.aid) .. ' subtitles=' .. tostring(saved.sid))
    end
    restoring = false
    ready = true
end)
for _, property in ipairs({'aid', 'sid', 'audio-delay', 'sub-delay', 'track-list'}) do
    mp.observe_property(property, 'native', persist)
end
