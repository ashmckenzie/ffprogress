#!/usr/bin/env ruby

STDOUT.sync = true

TOTAL = (ARGV[0] || 33).to_i

def nice_time(i)
  h = i / 3600
  m = (i % 3600) / 60
  s = i % 60

  '%02d:%02d:%02d' % [ h, m, s ]
end

puts "Duration: %s" % [ nice_time(TOTAL) ]

duration = 1

while duration <= TOTAL do
  print "\rtime=%s" % [ nice_time(duration) ]
  break if duration >= TOTAL
  wait = rand(2) + 1
  duration += wait
  sleep(wait)
end
