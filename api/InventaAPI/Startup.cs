using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.OpenApi.Models;
using Service;
using Microsoft.Identity.Web;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.Authorization;

namespace DFDSServiceAPI
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddAutoMapper(typeof(Startup));

            services.AddCors(options =>
                options.AddPolicy("GlobalPolicy",
                    builder => builder.AllowAnyOrigin().AllowAnyHeader().AllowAnyMethod()));

            services.AddControllers();
            services.AddSwaggerGen(c =>
            {
                if (bool.Parse(Configuration.GetSection("INVENTA_API_AUTH_ENABLE").Value))
                {
                    c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
                    {
                        In = ParameterLocation.Header,
                        Description = "Please insert JWT with Bearer into field",
                        Name = "Authorization",
                        Type = SecuritySchemeType.ApiKey
                    });

                    c.AddSecurityRequirement(new OpenApiSecurityRequirement {
                    {
                        new OpenApiSecurityScheme
                        {
                        Reference = new OpenApiReference
                        {
                            Type = ReferenceType.SecurityScheme,
                            Id = "Bearer"
                        }
                        },
                        new string[] { }
                    }
                    });
                }

                c.SwaggerDoc("v1", new OpenApiInfo { Title = "InventaAPI", Version = "v1" });
            });

            services.AddServiceProxyServiceCollection(Configuration);
            
            if (bool.Parse(Configuration.GetSection("INVENTA_API_AUTH_ENABLE").Value))
            {
                ConfigureAuth(services);
            }
            else
            {
                services.AddMvcCore(opts =>
                {
                    opts.EnableEndpointRouting = false;
                    opts.Filters.Add(new AllowAnonymousFilter());
                });
            }
        }

        protected virtual void ConfigureAuth(IServiceCollection services)
        {
            services.AddMvcCore(options => options.EnableEndpointRouting = false)
                    .AddAuthorization();
            services.AddMicrosoftIdentityWebApiAuthentication(Configuration, "AzureAD");

            var scopeRequirementPolicy = new AuthorizationPolicyBuilder().RequireAuthenticatedUser().Build();

            services.Configure<MvcOptions>(options =>
                    options.Filters.Add(new AuthorizeFilter(scopeRequirementPolicy))
                );

        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseSwagger();
            app.UseSwaggerUI(c =>
            {
                c.SwaggerEndpoint("/swagger/v1/swagger.json", "Inventa API v1");
                c.ConfigObject.AdditionalItems.Add("syntaxHighlight", false);
            });

            app.UseCors("GlobalPolicy");

            app.UseRouting();

            if (bool.Parse(Configuration.GetSection("INVENTA_API_AUTH_ENABLE").Value))
            {
                app.UseAuthentication();
            }

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });

            app.UseMvc();
        }
    }
}
