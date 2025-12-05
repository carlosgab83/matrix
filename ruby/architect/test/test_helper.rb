# typed: false
# frozen_string_literal: true

$LOAD_PATH.unshift File.expand_path('../lib', __dir__)

require 'minitest/autorun'
require 'minitest/reporters'

# Configure colored output
Minitest::Reporters.use! [
  Minitest::Reporters::SpecReporter.new
]

# Silence logs during tests
require 'architect/shared/app_logger'
Architect::Shared::AppLogger.setup(log_level: 'error')

require 'architect/domain/rule'
require 'architect/domain/action'
require 'architect/domain/condition'
require 'architect/domain/registry'
require 'architect/domain/rule_builder'
require 'architect/domain/rule_evaluator'
require 'architect/domain/errors/invalid_rule_definition'
