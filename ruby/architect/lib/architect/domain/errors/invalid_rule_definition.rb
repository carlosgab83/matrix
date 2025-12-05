# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'

module Architect
  module Domain
    module Errors
      class InvalidRuleDefinition < StandardError
        extend T::Sig
      end
    end
  end
end
