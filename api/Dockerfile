FROM mcr.microsoft.com/dotnet/sdk:5.0 AS build
WORKDIR /src
COPY . .
RUN dotnet restore "InventaAPI/InventaAPI.csproj"
RUN dotnet publish -c Release -o /app/out "InventaAPI/InventaAPI.csproj"

FROM mcr.microsoft.com/dotnet/aspnet:5.0 AS runtime
WORKDIR /app
EXPOSE 8080

RUN groupadd -r appuser
RUN useradd -r -g appuser appuser
RUN chown -R appuser:appuser /app
ENV ASPNETCORE_URLS http://+:8080
USER appuser

WORKDIR /app
COPY --from=build /app/out .
ENTRYPOINT ["dotnet", "InventaAPI.dll"]