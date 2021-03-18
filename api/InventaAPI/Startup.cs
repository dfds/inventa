using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.OpenApi.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Service;
using AutoMapper;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.AzureAD.UI;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.Identity.Web;
using Microsoft.Identity.Web.Resource;
using AadIssuerValidator = DFDSServiceAPI.Authentication.AadIssuerValidator;

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

                c.SwaggerDoc("v1", new OpenApiInfo { Title = "InventaAPI", Version = "v1" });
            });

            services.AddServiceProxyServiceCollection(Configuration);
            
            ConfigureAuth(services);
        }

        protected virtual void ConfigureAuth(IServiceCollection services)
        {
            services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = "AzureADBearer";
                options.DefaultChallengeScheme = "AzureADBearer";
            }).AddAzureADBearer(options =>
            {
              Configuration.Bind("AzureAd", options);  
            });

            services.Configure<JwtBearerOptions>(AzureADDefaults.JwtBearerAuthenticationScheme, options =>
            {
                options.Authority += "/v2.0";
                options.TokenValidationParameters.ValidAudiences = new string[] { options.Audience, $"api://{options.Audience}" };
                options.TokenValidationParameters.IssuerValidator = AadIssuerValidator.ValidateAadIssuer;
            });
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseSwagger();
            app.UseSwaggerUI(c => c.SwaggerEndpoint("/swagger/v1/swagger.json", "Inventa API v1"));

            app.UseCors("GlobalPolicy");

            app.UseRouting();

            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });
        }
    }
}
