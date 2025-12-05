# typed: false
# frozen_string_literal: true

# Disable stdout buffering for Docker logs
$stdout.sync = true
$stderr.sync = true

$LOAD_PATH.unshift(File.expand_path('lib', __dir__))

require 'architect/shared/app_logger'
require 'architect/domain/registry'
require 'architect/handler/app'

# Initialize logger with log level from environment
log_level = ENV['LOG_LEVEL'] || 'info'
Architect::Shared::AppLogger.setup(log_level: log_level)

rules_path = File.expand_path('lib/architect/rules', __dir__)
Architect::Domain::Registry.reload_from_dir!(rules_path)

puts Architect::Domain::Registry.rules.inspect

Architect::Handler::App.run
