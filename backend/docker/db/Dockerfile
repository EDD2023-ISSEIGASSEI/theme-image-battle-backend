FROM --platform=linux/x86_64 mysql:8.0-debian

#MySQL設定ファイルをイメージ内にコピー
COPY ./docker/db/my.cnf /etc/mysql/conf.d/my.cnf

# CLIで日本語設定
RUN apt-get update \
  && apt-get install -y locales \
  && sed -i -E 's/# (ja_JP.UTF-8)/\1/' /etc/locale.gen \
  && locale-gen \
  && update-locale LANG=ja_JP.UTF-8 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
ENV LC_ALL ja_JP.UTF-8

CMD ["mysqld"]