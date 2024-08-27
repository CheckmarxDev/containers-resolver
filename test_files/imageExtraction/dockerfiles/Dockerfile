# FROM ddd:test
FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build
# Install ReSharper GlobalTools
RUN dotnet tool install --global JetBrains.ReSharper.GlobalTools --version 2021.3.0

# Restore nuget packages
WORKDIR /src
COPY /src/*.sln /src/*/*.csproj /src/Nuget.config ./
RUN for file in $(ls *.csproj); do mkdir -p ${file%.*}/ && cp $file ${file%.*}/; done
RUN dotnet restore *.sln --configfile Nuget.config -nowarn:msb3202,nu1503

# Build
COPY /src/ .
RUN dotnet build *.sln


FROM build AS test
# Run ReSharper GlobalTools
RUN /root/.dotnet/tools/jb inspectcode /src/*.sln -o=output.txt -s=WARNING --verbosity=WARN --no-build -f=Text -a\
 && grep -E " {5}" output.txt | tee issues.txt && if [ -s issues.txt ]; then echo "Please fix the above warnings." && exit 1; fi

FROM build AS publish
RUN dotnet publish -c Release -o /app Bootstrap/Bootstrap.csproj


FROM mcr.microsoft.com/dotnet/aspnet:6.0
RUN apt-get update
RUN apt-get -y install wget
RUN apt-get -y install netcat

WORKDIR /app
COPY --from=publish /app .

EXPOSE 80
ENV ASPNETCORE_URLS http://*:80

CMD ["dotnet", "app.Bootstrap.dll"]
