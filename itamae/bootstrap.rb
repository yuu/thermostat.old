module RecipeHelper
    def include_cookbook(name)
        path = File.expand_path("../../../cookbook/#{name}", @recipe.path)
        include_recipe path
    end
end

MItamae::RecipeContext.send(:include, RecipeHelper)

include_recipe File.join("roles", node[:role], "default.rb")
