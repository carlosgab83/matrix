# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'logger'

module Architect
  module Shared
    class AppLogger
      class << self
        extend T::Sig

        sig { params(log_level: String).returns(Logger) }
        def setup(log_level: 'info')
          logger = Logger.new($stdout)
          logger.level = parse_level(log_level)
          logger.formatter = proc do |severity, datetime, _progname, msg|
            "[#{datetime.strftime('%Y-%m-%d %H:%M:%S')}] #{severity.ljust(5)} -- #{msg}\n"
          end
          @instance = T.let(logger, T.nilable(Logger))
          logger
        end

        sig { returns(Logger) }
        def instance
          @instance = T.let(nil, T.nilable(Logger)) unless defined?(@instance)
          @instance ||= setup
        end

        sig { params(message: String).void }
        def debug(message)
          instance.debug(message)
        end

        sig { params(message: String).void }
        def info(message)
          instance.info(message)
        end

        sig { params(message: String).void }
        def warn(message)
          instance.warn(message)
        end

        sig { params(message: String).void }
        def error(message)
          instance.error(message)
        end

        sig { params(message: String).void }
        def fatal(message)
          instance.fatal(message)
        end

        private

        sig { params(level: String).returns(Integer) }
        def parse_level(level)
          case level.to_s.downcase
          when 'debug' then Logger::DEBUG
          when 'warn', 'warning' then Logger::WARN
          when 'error' then Logger::ERROR
          when 'fatal' then Logger::FATAL
          else Logger::INFO
          end
        end
      end
    end
  end
end
